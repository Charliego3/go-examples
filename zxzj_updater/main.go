package main

type Updater struct {
	Name string
	URL  string
}

type UpdateFunc interface {
	Update()
}

var updates = []UpdateFunc{
	FalconAndWinterSoldier{
		Updater{
			Name: "猎鹰与冬兵",
			URL:  "https://www.zxzj.me/video/3212-1-5.html",
		},
	},
}

func main() {
	for _, updater := range updates {
		updater.Update()
	}
}
