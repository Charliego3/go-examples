package main

import (
	"github.com/caseymrm/menuet"
	"strconv"
)

const (
	number   = "1234567890"
	azs      = "abcdefghijklmnopqrstuvwxyz"
	AZs      = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specials = "-,./?><;'[]`!@#$%^&*()"
)

func gen() {

}

func generatePassword() menuet.MenuItem {
	item := menuet.MenuItem{
		Text:    "Generate password",
		Clicked: gen,
		Children: func() []menuet.MenuItem {
			return []menuet.MenuItem{
				numbers(),
				az(),
				AZ(),
				special(),
				length(),
			}
		},
	}
	return item
}

func length() menuet.MenuItem {
	const key = "password_length"
	pwdLength := menuet.Defaults().String(key)
	if pwdLength == "" {
		pwdLength = "click to set length, or use 15"
	}
	return menuet.MenuItem{
		Text: pwdLength,
		Clicked: func() {
			r := menuet.App().Alert(menuet.Alert{
				MessageText:     "Please input password length",
				InformativeText: "this is using password generate length",
				Buttons:         []string{"OK", "Cancel"},
				Inputs:          []string{""},
			})
			if r.Button != 0 {
				return
			}

			length := r.Inputs[0]
			if length == "" {
				notification("Set password length error", "you input is empty", "password length can not be empty")
				return
			}

			_, err := strconv.Atoi(length)
			if err != nil {
				notification("Set password length error", "you input is not a number", "password length only can input number, you can try again.")
				return
			}

			menuet.Defaults().SetString(key, "Password length: "+length)
		},
	}
}

func special() menuet.MenuItem {
	const key = "password_special"
	state := menuet.Defaults().Boolean(key)
	return menuet.MenuItem{
		Text:  "Special Symbols",
		State: state,
		Clicked: func() {
			menuet.Defaults().SetBoolean(key, !state)
		},
	}
}

func AZ() menuet.MenuItem {
	const key = "password_AZ"
	state := menuet.Defaults().Boolean(key)
	return menuet.MenuItem{
		Text:  "A ~ Z",
		State: state,
		Clicked: func() {
			menuet.Defaults().SetBoolean(key, !state)
		},
	}
}

func az() menuet.MenuItem {
	const key = "password_az"
	state := menuet.Defaults().Boolean(key)
	return menuet.MenuItem{
		Text:  "a ~ z",
		State: state,
		Clicked: func() {
			menuet.Defaults().SetBoolean(key, !state)
		},
	}
}

func numbers() menuet.MenuItem {
	const key = "password_09"
	state := menuet.Defaults().Boolean(key)
	return menuet.MenuItem{
		Text:  "0 ~ 9",
		State: state,
		Clicked: func() {
			menuet.Defaults().SetBoolean(key, !state)
		},
	}
}
