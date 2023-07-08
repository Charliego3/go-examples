package main

import (
	"github.com/charliego93/logger"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

func main() {
	brower := rod.New().MustConnect()
	defer brower.MustClose()

	numberSelector := "input[placeholder=手机号]"
	passwordSelector := "input[type=password]"

	page := brower.MustPage("https://www.bw6.com/cn/register").MustWaitLoad()
	page.MustElement(numberSelector).MustInput("18929345654").MustType(input.Enter)
	page.MustElement(passwordSelector).MustInput("ZXCvbn,./123").MustType(input.Enter)

	number := page.MustElement(numberSelector).MustText()
	logger.Info("have enter phone number", "value", number)
	password := page.MustElement(passwordSelector).MustText()
	logger.Info("have enter password", "value", password)

	page.MustElement("button").MustClick()

	text, err := page.HTML()
	if err != nil {
		panic(err)
	}
	logger.Print(text)
}
