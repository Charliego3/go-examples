package main

import (
	"context"
	"github.com/transerver/commons/logger"
	"testing"
	"time"
)

func TestWebsocket(t *testing.T) {
	ctx := context.Background()
	receiver := NewReceiver("push_user_record")
	depthReceiver := NewReceiver(QuickDepthType)
	dialer := NewWebsocketDialer(ctx, "ws://ttapi2.100-130.net/websocket", receiver, depthReceiver)
	user := User{
		ID:        362652,
		Username:  "15200000021",
		APIKey:    "c5c8e7a7-92e2-4314-9ca3-c869911d7905",
		APISecret: "add4417a-2554-4d4a-9832-8ce1cc30c6aa",
	}
	err := dialer.Connect(user)
	if err != nil {
		t.Fatal(err)
	}

	message := "{\"accesskey\":\"%s\",\"channel\":\"push_user_record\",\"event\":\"%s\",\"market\":\"btcqcdefault\""
	sign := dialer.sign(message, true)
	err = dialer.SendMessage(sign)
	if err != nil {
		t.Fatal(err)
	}

	err = dialer.SubscribeQuickDepth("btcqc")
	if err != nil {
		t.Fatal(err)
	}

	prevTime := time.Time{}

	for {
		select {
		case msg := <-receiver.C:
			logger.Warnf("收到委托: %s", msg)
		case msg := <-depthReceiver.C:
			now := time.Now()
			if !prevTime.IsZero() && now.Sub(prevTime) < time.Minute {
				continue
			}
			logger.Debugf("收到深度: %s", msg)
			prevTime = now
		}
	}
}
