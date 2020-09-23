package main

import (
	"encoding/json"
	"github.com/whimthen/kits/logger"
	"net/http"
)

func DoBodyPost(url string, body map[string]string, cookies map[string]string) {
	bytes, err := json.Marshal(body)
	if err != nil {
		logger.Error("Post Error, can't marshal body to bytes, %+v", err)
		return
	}

	logger.Debug("Post body is %s", bytes)
	logger.Info("Method: %s", http.MethodPost)
	// http.NewRequest(http.MethodPost, url)
}
