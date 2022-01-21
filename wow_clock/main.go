package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kataras/golog"
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
	"github.com/robfig/cron/v3"
	"github.com/whimthen/temp/macdriver-gui/widgets/alert"
	"github.com/whimthen/temp/macdriver-gui/widgets/statusBar"
	"net/http"
	"time"
)

type clock struct {
	opened bool
	cronId cron.EntryID
}

func (c *clock) Title() string {
	if c.opened {
		return "✅ 定时打卡"
	}
	return "定时打卡"
}

type chmsg struct {
	succ bool
	msg  string
}

var (
	c      = &clock{opened: true}
	logger = cron.PrintfLogger(golog.Default)
	cc     = cron.New(cron.WithChain(cron.Recover(logger)), cron.WithLogger(logger), cron.WithSeconds())
	client = resty.New()
	ch     = make(chan chmsg)
)

func init() {
	golog.SetTimeFormat("2006-01-02 15:04:05.000000")
	client.SetRetryWaitTime(time.Minute * 2)
	client.SetTimeout(2 * time.Minute)
	client.SetHeader("Content-Type", "application/json; charset=utf-8")
	client.SetHeader("host", "approve.yowoworld.cc")
	client.AddRetryCondition(func(response *resty.Response, err error) bool {
		return err != nil || response.StatusCode() != http.StatusOK
	})
	addCron()
	cc.Start()
}

func getTitle(t time.Time) string {
	if t.After(time.Date(t.Year(), t.Month(), t.Day(), 19, 0, 0, 0, time.Local)) ||
		t.Before(time.Date(t.Year(), t.Month(), t.Day(), 9, 0, 0, 0, time.Local)) {
		return "🕘 Wow"
	} else {
		return "⏰ Wow"
	}
}

func business(item cocoa.NSStatusItem) {
	go func() {
		for {
			select {
			case msg := <-ch:
				now := time.Now()
				hour := now.Hour()
				message := "Wow " + now.Format("2006-01-02 15:04:05.000") + " \n"
				if hour > 12 {
					hour = 9
					message += "下班"
				} else {
					hour = 19
					message += "上班"
				}
				core.Dispatch(func() {
					item.Button().SetTitle(getTitle(now))
					showMsg(message + msg.msg)
				})
			}
		}
	}()
}

func main() {
	app := statusBar.NewStatusBarApp(getTitle(time.Now()), business)
	item := cocoa.NSMenuItem_New()
	item.SetTitle(c.Title())
	item.SetAction(objc.Sel("clock:"))
	cocoa.DefaultDelegateClass.AddMethod("clock:", clockIn(item))
	app.Menu.AddItem(item)

	app.AddItemSeparator()
	app.AddTerminateItem()

	app.Run()
}

func clockIn(item cocoa.NSMenuItem) func(objc.Object) {
	return func(_ objc.Object) {
		if !c.opened { // 打开
			addCron()
			c.opened = true
			item.SetTitle(c.Title())
			showMsg("开启定时打卡成功")
		} else {
			cc.Remove(c.cronId)
			c.opened = false
			item.SetTitle(c.Title())
			showMsg("关闭定时打卡成功")
		}
	}
}

func addCron() {
	id, err := cc.AddFunc("0 0 9,19 * * 1,2,3,4,5,6", func() {
		//id, err := cc.AddFunc("0 0/1 * * * ?", func() {
		hour := time.Now().Hour()
		if hour > 12 {
			request("144", "1637665332707", "71b4dff1679e220c956527a358b3257f",
				`{"wowId": "82a0f47a80e3431324e0efac32036d36", "clockInUserName": Whimthen, "clockInAddress": %E6%B7%B1%E5%9C%B3, "clockInType": "0", "remark": }`)
		} else {
			request("144", "1637803346765", "c7415cbd0d2cea64949993504a980b91",
				`{"wowId": "82a0f47a80e3431324e0efac32036d36", "clockInUserName": Whimthen, "clockInAddress": %E6%B7%B1%E5%9C%B3, "clockInType": "1", "remark": }`)
		}
		//golog.Info(open)
	})
	if err != nil {
		panic(err)
	}
	c.cronId = id
	golog.Infof("Time: %s, Id: %d", time.Now().String(), id)
}

func request(length, timestamp, sign, body string) {
	response, err := client.R().SetHeader("Content-Length", length).SetHeader("timestamp", timestamp).SetHeader("sign", sign).SetBody(body).
		Post("https://approve.yowoworld.cc/dingteam/AttendanceController/clockIn")
	if err != nil {
		panic(err)
	}

	code := response.StatusCode()
	if code == http.StatusOK {
		ch <- chmsg{
			succ: true,
			msg:  fmt.Sprintf("打卡成功: %s", response.Body()),
		}
	} else {
		ch <- chmsg{
			succ: false,
			msg:  fmt.Sprintf("打卡失败: %s", response.Body()),
		}
	}
}

func showMsg(msg string) {
	nsAlert := alert.NewNSAlert()
	nsAlert.SetAlertStyle(alert.Informational)
	nsAlert.SetMessageText(msg)
	nsAlert.AddButtonWithTitle("OK")
	nsAlert.Show()
}
