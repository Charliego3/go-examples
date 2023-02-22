package websocket

import (
	"fmt"
	"github.com/transerver/commons/utils"
)

type RequestMsg struct {
	*JsonMessage
	Event     string  `json:"event,omitempty"`
	Channel   string  `json:"channel,omitempty"`
	AccessKey string  `json:"accesskey,omitempty"`
	Market    *string `json:"market,omitempty"`
	Sign      *string `json:"sign,omitempty"`
}

func (req *RequestMsg) Signed(secretKey string) IMessage {
	placeholder := `"accesskey":"%s","channel":"%s","event":"%s"`
	if req.Market != nil {
		*req.Market += "default"
		placeholder += `,"market":"` + *req.Market + `"`
	}

	sha := HexSha1([]byte(secretKey))
	placeholder = fmt.Sprintf(placeholder, req.AccessKey, req.Channel, req.Event)
	sign := HmacMD5(utils.Bytes("{"+placeholder+"}"), sha)
	req.Sign = &sign
	return req
}
