package main

import (
	"github.com/kataras/golog"
	"github.com/muesli/termenv"
)

func main() {
	golog.Info(termenv.HasDarkBackground())
}
