package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/transerver/commons/utils"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	json "github.com/json-iterator/go"
	"github.com/transerver/commons/logger"
)

type DataType string

const (
	QuickDepthType     DataType = "quickDepth"
	UserAssetType               = "userAsset"
	UserIncrAssetType           = "userIncrAsset"
	UserIncrRecordType          = "userIncrRecord"

	quickDepthMessage     = "{\"event\":\"%s\",\"channel\":\"%s_quick_depth\",\"length\":\"20\"}"
	userIncrRecordMessage = "{\"accesskey\":\"%s\",\"channel\":\"push_user_incr_record\",\"event\":\"%s\",\"market\":\"%sdefault\""
	userIncrAssetMessage  = "{\"accesskey\":\"%s\",\"channel\":\"push_user_incr_asset\",\"event\":\"%s\""
	userAssetMessage      = "{\"accesskey\":\"%s\",\"channel\":\"push_user_asset\",\"event\":\"%s\""
	unSubUserAssetMessage = "{\"accesskey\":\"%s\",\"channel\":\"push_user_asset\",\"event\":\"%s\""
)

var (
	ping = []byte("ping")
	pong = []byte("pong")
)

type WebsocketDialer struct {
	conn     *websocket.Conn
	address  string
	Logger   *logger.Logger
	receives map[DataType]*Receiver
	done     bool
	mutex    sync.Mutex
	rmutex   sync.RWMutex
	ctx      context.Context

	user User
}

type Receiver struct {
	T DataType
	C chan []byte
}

func NewReceiver(t DataType) *Receiver {
	return &Receiver{
		C: make(chan []byte),
		T: t,
	}
}

func NewWebsocketDialer(ctx context.Context, address string, rs ...*Receiver) *WebsocketDialer {
	return NewWebsocketDialerWithSuffix(ctx, address, "", rs...)
}

func NewWebsocketDialerWithSuffix(ctx context.Context, address, suffix string, rs ...*Receiver) *WebsocketDialer {
	if utils.NotBlank(suffix) {
		suffix = ":" + suffix
	}

	w := &WebsocketDialer{
		address:  address,
		Logger:   logger.NewLogger(logger.WithPrefix(address + suffix)),
		receives: make(map[DataType]*Receiver),
		ctx:      ctx,
	}

	for _, r := range rs {
		w.SetReceiver(r)
	}

	return w
}

func (w *WebsocketDialer) Receiver(dataType DataType) (r *Receiver, ok bool) {
	w.rmutex.RLock()
	defer w.rmutex.RUnlock()

	r, ok = w.receives[dataType]
	return
}

func (w *WebsocketDialer) SetReceiver(r *Receiver) {
	w.rmutex.Lock()
	defer w.rmutex.Unlock()

	w.receives[r.T] = r
}

func (w *WebsocketDialer) RemoveReceiver(dataType DataType) {
	if r, ok := w.receives[dataType]; ok {
		close(r.C)
	}

	w.rmutex.Lock()
	defer w.rmutex.Unlock()

	delete(w.receives, dataType)
}

func (w *WebsocketDialer) SubscribeQuickDepth(market string) error {
	message := fmt.Sprintf(quickDepthMessage, getEvent(true), market)
	return w.SendMessage(message)
}

func (w *WebsocketDialer) SubscribeUserIncrAsset() error {
	message := w.sign(userIncrAssetMessage, true)
	return w.SendMessage(message)
}

func (w *WebsocketDialer) SubscribeUserAsset() error {
	message := w.sign(userAssetMessage, true)
	return w.SendMessage(message)
}

func (w *WebsocketDialer) UnSubscribeUserAsset() error {
	message := w.sign(unSubUserAssetMessage, false)
	err := w.SendMessage(message)
	w.RemoveReceiver(UserAssetType)
	return err
}

func (w *WebsocketDialer) SubscribeIncrRecord(market Market) error {
	message := w.sign(userIncrRecordMessage, true, market.Name)
	return w.SendMessage(message)
}

func (w *WebsocketDialer) sign(message string, sub bool, args ...any) string {
	params := []any{w.user.APIKey, getEvent(sub)}

	if len(args) > 0 {
		params = append(params, args...)
	}

	message = fmt.Sprintf(message+"}", params...)
	sign := HmacMD5(message, w.user.APISecret)
	return message[:len(message)-1] + ",\"sign\":\"" + sign + "\"}"
}

func (w *WebsocketDialer) Connect(us ...User) error {
	if len(us) > 0 {
		w.user = us[0]
	}

	conn, _, err := websocket.DefaultDialer.Dial(w.address, nil)
	if err != nil {
		w.Logger.Errorf("websocket连接出错: %v", err)
		return err
	}

	w.conn = conn
	go w.accept()
	go w.heartBeat()
	return nil
}

func (w *WebsocketDialer) Shutdown() error {
	defer func() {
		w.rmutex.RLock()
		for _, r := range w.receives {
			close(r.C)
		}
		w.rmutex.RUnlock()
		w.Logger.Debug("已关闭websocket连接")
	}()

	message := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	w.mutex.Lock()
	w.done = true
	err := w.conn.WriteMessage(websocket.CloseMessage, message)
	w.mutex.Unlock()
	if err != nil {
		return err
	}
	err = w.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (w *WebsocketDialer) SendMessage(message string, t ...int) error {
	msgType := websocket.TextMessage
	if len(t) > 0 {
		msgType = t[0]
	}

	w.mutex.Lock()
	err := w.conn.WriteMessage(msgType, []byte(message))
	w.mutex.Unlock()
	if err != nil {
		w.Logger.Errorf("发送消息出错: %v, Message: %s, MsgType: %d", err, message, msgType)
	}
	w.Logger.Infof("发送消息: %s", message)
	return err
}

func (w *WebsocketDialer) accept() {
	defer func() {
		if err := recover(); err != nil {
			w.Logger.Errorf("Websocket处理消息出错: %v", err)
		}
	}()

	for {
		_, message, err := w.conn.ReadMessage()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}

			w.Logger.Errorf("读取消息出错: %v", err)
			continue
		}

		if bytes.Compare(message, pong) == 0 {
			w.Logger.Debugf("收到Pong消息: %s", message)
			continue
		}

		dataType := DataType(json.Get(message, "dataType").ToString())

		if dataType != QuickDepthType {
			logger.Errorf("消息: %s", message)
		}

		r, ok := w.Receiver(dataType)
		if !ok {
			w.Logger.Warnf("消息没有接收者: %s, Message: %s", dataType, message)
			continue
		}

		r.C <- message
	}
}

func (w *WebsocketDialer) heartBeat() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.mutex.Lock()
			if w.done {
				w.mutex.Unlock()
				return
			}

			err := w.conn.WriteMessage(websocket.TextMessage, ping)
			w.mutex.Unlock()
			if err != nil {
				w.Logger.Errorf("发送ping消息失败: %v", err)
			}
		case <-w.ctx.Done():
			w.Logger.Debug("即将关闭websocket链接")
			err := w.Shutdown()
			if err != nil {
				w.Logger.Error("关闭websocket链接出错:", err)
			}
			return
		}
	}
}

func getEvent(sub bool) string {
	if sub {
		return "addChannel"
	}
	return "removeChannel"
}

// HmacMD5 sign the websocket message
func HmacMD5(message, secretKey string) (hmacSign string) {
	h := hmac.New(md5.New, []byte(digest(secretKey)))
	h.Write([]byte(message))
	hmacSign = hex.EncodeToString(h.Sum(nil))
	return
}

// SHA1 加密
func digest(secretKey string) (digest string) {
	hash := sha1.New()
	hash.Write([]byte(secretKey))
	digest = hex.EncodeToString(hash.Sum(nil))
	return
}
