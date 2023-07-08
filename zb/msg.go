package zb

import (
	"fmt"

	"github.com/charliego93/websocket"
	"github.com/transerver/commons/utils"
)

type RequestMsg struct {
	*websocket.JsonMessage
	Event     string  `json:"event,omitempty"`
	Channel   string  `json:"channel,omitempty"`
	AccessKey string  `json:"accesskey,omitempty"`
	Market    *string `json:"market,omitempty"`
	Sign      *string `json:"sign,omitempty"`
}

func (req *RequestMsg) Signed(secretKey string) websocket.IMessage {
	placeholder := `"accesskey":"%s","channel":"%s","event":"%s"`
	if req.Market != nil {
		*req.Market += "default"
		placeholder += `,"market":"` + *req.Market + `"`
	}

	sha := websocket.HexSha1([]byte(secretKey))
	placeholder = fmt.Sprintf(placeholder, req.AccessKey, req.Channel, req.Event)
	sign := websocket.HmacMD5(utils.Bytes("{"+placeholder+"}"), sha)
	req.Sign = &sign
	return req
}
