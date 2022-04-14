package main

import (
	"context"
	json "github.com/json-iterator/go"
	"go.uber.org/atomic"
	"strings"
)

var currentDepth atomic.Value

func SubscribeQuickDepth(ctx context.Context, market Market) error {
	dr := NewReceiver(QuickDepthType)
	dialer := NewWebsocketDialer(ctx, Settings.Websocket.Address+"/"+strings.ToLower(market.Symbol), dr)

	err := dialer.Connect()
	if err != nil {
		return err
	}

	err = dialer.SubscribeQuickDepth(market.Name)
	if err != nil {
		return err
	}

	go acceptDepth(ctx, dr, dialer)

	return nil
}

func Depth() QuickDepth {
	depth := currentDepth.Load()
	if depth == nil {
		return QuickDepth{}
	}

	return *depth.(*QuickDepth)
}

func acceptDepth(ctx context.Context, dr *Receiver, dialer *WebsocketDialer) {
	for {
		select {
		case m, ok := <-dr.C:
			if !ok {
				dr.C = nil
				break
			}

			var depth QuickDepth
			err := json.Unmarshal(m, &depth)
			if err != nil {
				break
			}

			currentDepth.Store(&depth)
		case <-ctx.Done():
			dialer.Logger.Debug("退出")
			currentDepth.Store(&QuickDepth{})
			return
		}
	}
}
