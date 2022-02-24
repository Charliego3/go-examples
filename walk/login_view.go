package main

import (
	"fmt"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

type LoginWindow struct {
	*walk.Dialog

	Number   string
	Username string
	Password string
	Remember bool
}

const (
	number   = "number"
	username = "username"
	password = "password"
	remember = "remember"
)

func NewLoginWindow() (*LoginWindow, error) {
	formVSpacer := declarative.VSpacer{
		ColumnSpan: 2,
		Size:       5,
	}
	formFont := NewFont(15, true)

	var db *walk.DataBinder
	loginWindow := new(LoginWindow)
	rememberVal, _ := settings.Get(remember)
	if rememberVal == "true" {
		loginWindow.Remember = true
		loginWindow.Number, _ = settings.Get(number)
		loginWindow.Username, _ = settings.Get(username)
		loginWindow.Password, _ = settings.Get(password)
	}

	err := declarative.Dialog{
		AssignTo:  &loginWindow.Dialog,
		Title:     WindowTitle,
		Layout:    declarative.Grid{MarginsZero: true, Columns: 2},
		FixedSize: true,
		Children: []declarative.Widget{
			declarative.ImageView{
				Background: declarative.SolidColorBrush{Color: walk.RGB(255, 255, 255)},
				Image:      "resources/logo.jpg",
				Mode:       declarative.ImageViewModeCenter,
				MinSize:    declarative.Size{Width: 400},
			},
			declarative.Composite{
				DataBinder: declarative.DataBinder{
					AssignTo:       &db,
					Name:           "loginModel",
					DataSource:     loginWindow,
					ErrorPresenter: declarative.ToolTipErrorPresenter{},
				},
				MinSize: declarative.Size{Width: 300, Height: 600},
				Layout:  declarative.Grid{Columns: 2},
				Children: []declarative.Widget{
					declarative.VSpacer{
						ColumnSpan: 2,
					},

					declarative.Label{
						Text:          "黄土塬餐饮",
						ColumnSpan:    2,
						TextAlignment: declarative.AlignCenter,
						Font:          NewFont(30, true),
					},

					declarative.VSpacer{
						ColumnSpan: 2,
						Size:       30,
					},

					declarative.Label{
						Text: "商户编号:",
						Font: formFont,
					},
					declarative.LineEdit{
						Text:               declarative.Bind("Number"),
						AlwaysConsumeSpace: false,
						Font:               formFont,
					},

					formVSpacer,

					declarative.Label{
						Text:          "用户名:",
						TextAlignment: declarative.AlignFar,
						Font:          formFont,
					},
					declarative.LineEdit{
						Text: declarative.Bind("Username"),
						Font: formFont,
					},

					formVSpacer,

					declarative.Label{
						Text:          "密码:",
						TextAlignment: declarative.AlignFar,
						Font:          formFont,
					},
					declarative.LineEdit{
						Text:         declarative.Bind("Password"),
						Font:         formFont,
						PasswordMode: true,
					},

					formVSpacer,

					declarative.Label{
						Text: "记住密码:",
						Font: formFont,
					},
					declarative.CheckBox{
						Checked: declarative.Bind("Remember"),
						Text:    "勾选后下次登录不用手动输入密码",
						Font:    NewFont(10),
					},

					declarative.VSpacer{
						ColumnSpan: 2,
						Size:       30,
					},

					declarative.PushButton{
						ColumnSpan: 2,
						Text:       "登录",
						Font:       NewFont(20, true),
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								walk.MsgBox(loginWindow, "出错了", err.Error(), walk.MsgBoxIconError)
								return
							}

							_ = settings.Put(number, loginWindow.Number)
							_ = settings.Put(username, loginWindow.Username)
							_ = settings.Put(password, loginWindow.Password)
							_ = settings.Put(remember, fmt.Sprintf("%t", loginWindow.Remember))

							err := settings.Save()
							if err != nil {
								walk.MsgBox(loginWindow, "记住失败", fmt.Sprintf("%s", err.Error()), walk.MsgBoxIconError)
								return
							}

							if loginWindow.Number != "" {
								loginWindow.Accept()
							} else {
								walk.MsgBox(loginWindow, "登陆失败", "商户编号或用户名或密码错误", walk.MsgBoxIconError)
							}
						},
					},

					declarative.VSpacer{
						ColumnSpan: 2,
					},
				},
			},
		},
	}.Create(nil)

	if err == nil {
		Center(loginWindow, windowWidth, windowHeight)
	}

	return loginWindow, err
}
