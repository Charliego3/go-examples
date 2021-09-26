package main

import (
	"github.com/kataras/golog"
	"github.com/kyoto44/rain/magnet"
)

func main() {
	mag, err := magnet.New("magnet:?xt=urn:btih:5E1B7AFC69EF48F7D6792ACFD9DDD03FB6C9DA60")
	if err != nil {
		golog.Fatal(err)
	}

	golog.Error(mag.Name)
	golog.Error(mag.Peers)
	golog.Error(mag.Trackers)
	golog.Error(mag.InfoHash)
	golog.Error(mag.String())
}
