package main

import (
	"fmt"
	"github.com/kataras/golog"
)

type Args struct {
	ID      int64
	RobotId int64
	Env     string
	Market  string

	ProdStrategyDBURL string
	ProdEntrustDBURL  string
}

func parseArgs() bool {
	dbValid := false
	strategyDbValid := false
	entrustDbValid := false
	if args.ProdStrategyDBURL != "" {
		strategyDbValid = true
	}
	if args.ProdStrategyDBURL != "" {
		entrustDbValid = true
	}
	if strategyDbValid && !entrustDbValid {
		golog.Error("必须同时指定网格的数据库链接地址: user:password@tcp(ip:port)/dbname")
		return false
	} else if !strategyDbValid && entrustDbValid {
		golog.Error("必须同时指定盘口的数据库链接地址: user:password@tcp(ip:port)/dbname")
		return false
	} else if strategyDbValid && entrustDbValid {
		dbValid = true

		strategyDBURL = args.ProdStrategyDBURL
		entrustDBURL = args.ProdEntrustDBURL
		isProd = true
	}
	if !dbValid {
		hasContainsEnv := false
		for env := range Envs {
			if env == args.Env {
				hasContainsEnv = true
				break
			}
		}
		if !hasContainsEnv {
			golog.Error("必须指定环境, eg: --env[-e] 130")
			return false
		}

		env := Envs[args.Env]
		strategyDBURL = env + strategyDBName
	}

	idValid := false
	robotIdValid := false
	if args.ID > 0 {
		idValid = true
	}

	if args.RobotId > 0 {
		robotIdValid = true
	}

	if !idValid && !robotIdValid {
		golog.Error("网格记录ID[id]和机器人ID[robot]必须指定其一, eg: --id 1231 或 --robot 1231")
		return false
	}

	if idValid && !robotIdValid && args.Market == "" {
		golog.Error("必须指定市场, eg: --market btcqc")
		return false
	}

	if args.Market != "" {
		env := Envs[args.Env]
		entrustDBURL = env + fmt.Sprintf(entrustDBName, args.Market)
	}

	return true
}
