package main

import (
	"fmt"
	"github.com/guonaihong/gout"
	"github.com/shopspring/decimal"
	"github.com/transerver/commons/logger"
	"strings"
)

const (
	getUserInfoEndpoint   = "/api/fake/V1_0_0/getUserInfo"
	createUserKeyEndpoint = "/api/fake/V1_0_0/createUserKey"
	doRechargeEndpoint    = "/api/fake/V1_0_0/doRecharge"
	doRegisterEndpoint    = "/api/fake/V1_0_0/doRegister"
	getFundsEndpoint      = "/api/fake/V1_0_0/getFunds"

	password = "Zxcvbnm,./"
)

type UserInfoResp struct {
	UserId   int    `json:"userId,omitempty"`
	Username string `json:"username,omitempty"`
}

type UserKeyResp struct {
	APIKey    string `json:"apiKey"`
	APISecret string `json:"apiSecret"`
}

type UserFund struct {
	BTC UFund `json:"BTC"`
	QC  UFund `json:"QC"`
}

type UFund struct {
	Available decimal.Decimal `json:"available"`
	Freeze    decimal.Decimal `json:"freeze"`
	Total     decimal.Decimal `json:"total"`
}

func getRequestURL(endpoint string) string {
	return Settings.Basic.Domain + endpoint
}

func GetFunds(userId int, coin ...string) (UserFund, error) {
	var resp Response[UserFund]
	err := timeout.GET(getRequestURL(getFundsEndpoint)).SetQuery(gout.H{
		"userId": userId,
		"coins":  strings.Join(coin, ","),
	}).BindJSON(&resp).Do()
	if err != nil || !resp.Success() {
		logger.Errorf("获取资产失败: %d:[%s] -> %s%v", userId, strings.Join(coin, ","), resp.ResMsg.Message, err)
		return UserFund{}, err
	}

	return resp.Payload, nil
}

func Register(isTrading bool, phone ...int) (resp UserInfoResp, err error) {
	var username int
	if len(phone) == 0 {
		username = Settings.Basic.LastRegister
	} else {
		username = phone[0]
	}

	var response Response[UserInfoResp]
	err = timeout.GET(getRequestURL(doRegisterEndpoint)).SetQuery(gout.H{
		"userType": 1,
		"userName": username,
		"loginPwd": password,
	}).BindJSON(&response).Do()

	if err != nil {
		logger.Errorf("注册失败[%d], err: %v", username, err)
		return resp, err
	}

	if !response.Success() {
		logger.Errorf("注册失败[%d]: %s", username, response.ResMsg.Message)
		if username == Settings.Basic.LastRegister && strings.Contains(response.ResMsg.Message, "已注册") {
			Settings.Basic.LastRegister++
			return Register(isTrading, Settings.Basic.LastRegister)
		}
		return resp, fmt.Errorf("注册失败[%d]: %s", username, response.ResMsg.Message)
	}

	payload := response.Payload
	userKey, err := CreateUserKey(payload.UserId)
	if err != nil {
		logger.Errorf("创建用户APIkey失败, 重试, U: %d, N: %s", payload.UserId, payload.Username)
		Settings.Basic.LastRegister++
		return Register(isTrading, Settings.Basic.LastRegister)
		// return UserInfoResp{}, err
	}

	overrideUsersConfig(isTrading, User{
		ID:        payload.UserId,
		Username:  payload.Username,
		APIKey:    userKey.APIKey,
		APISecret: userKey.APISecret,
	})

	Settings.Basic.LastRegister++
	return payload, nil
}

func GetUserInfo(username int) (UserInfoResp, error) {
	var resp Response[UserInfoResp]
	err := timeout.GET(getRequestURL(getUserInfoEndpoint)).SetQuery(gout.H{"username": username}).BindJSON(&resp).Do()
	if err != nil {
		logger.Errorf("获取用户信息失败[%d]: %v", username, err)
		return UserInfoResp{}, err
	}

	if !resp.Success() {
		logger.Errorf("获取用户信息失败[%d]: %s", username, resp.ResMsg.Message)
		return UserInfoResp{}, fmt.Errorf("获取用户信息失败[%d]: %s", username, resp.ResMsg.Message)
	}

	logger.Debugf("用户信息: %v", resp)
	return resp.Payload, nil
}

func Recharge(userId int, coin string, amount decimal.Decimal) error {
	var resp Response[any]
	err := timeout.GET(getRequestURL(doRechargeEndpoint)).SetQuery(gout.H{
		"userId":   userId,
		"currency": coin,
		"amount":   amount,
	}).BindJSON(&resp).Do()
	if err != nil || !resp.Success() {
		logger.Errorf("充值失败[%d-%s-%s]: %v%s", userId, coin, amount.String(), err, resp.ResMsg.Message)
		return err
	}

	return nil
}

func CreateUserKey(userId int) (UserKeyResp, error) {
	var resp Response[UserKeyResp]
	err := timeout.GET(getRequestURL(createUserKeyEndpoint)).SetQuery(gout.H{"userId": userId}).BindJSON(&resp).Do()
	if err != nil || !resp.Success() {
		logger.Errorf("创建用户APIkey出错: %v%s", err, resp.ResMsg.Message)
		return UserKeyResp{}, err
	}
	return resp.Payload, nil
}

func overrideUsersConfig(isTrading bool, user User) {
	if isTrading {
		Settings.TradingUsers = append(Settings.TradingUsers, &user)
	} else {
		Settings.TradeUsers = append(Settings.TradeUsers, &user)
	}
	overrideConfig()
}
