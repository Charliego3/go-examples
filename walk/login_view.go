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

// func (lw *LoginWindow) Run() {
// 	(*lw.Dialog).Run()
// }

func NewLoginWindow() (*LoginWindow, error) {
	formVSpacer := declarative.VSpacer{
		ColumnSpan: 2,
		Size:       5,
	}
	formFont := NewFont(15, true)

	windowWidth := 800
	windowHeight := 600

	var db *walk.DataBinder
	loginWindow := new(LoginWindow)
	err := declarative.Dialog{
		AssignTo: &loginWindow.Dialog,
		Title:    WindowTitle,
		// Size:      declarative.Size{Width: windowWidth, Height: windowHeight},
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

							declarative.Dialog{
								Title:  "显示输入框的值",
								Layout: declarative.VBox{},
								Children: []declarative.Widget{
									declarative.Label{Text: loginWindow.Number},
									declarative.Label{Text: loginWindow.Username},
									declarative.Label{Text: loginWindow.Password},
									declarative.Label{Text: fmt.Sprintf("%t", loginWindow.Remember)},
									declarative.Label{Text: "Test dialog"},
								},
							}.Run(nil)

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
		width, height := GetWinScreen()
		walk.MsgBox(loginWindow, "屏幕长宽", fmt.Sprintf("长: %d, 宽: %d", width, height), walk.MsgBoxIconInformation)
		loginWindow.SetBounds(walk.Rectangle{
			X:      500,
			Y:      500,
			Width:  windowWidth,
			Height: windowHeight,
		})
	}

	return loginWindow, err
}
