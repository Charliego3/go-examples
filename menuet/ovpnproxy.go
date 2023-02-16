package main

import (
	"context"
	"fmt"
	"github.com/caseymrm/menuet"
	"github.com/go-vgo/robotgo/clipboard"
	"net/http"
	"strings"
)

var (
	opening = false
	server  *http.Server
)

const (
	listenAddress = "127.0.0.1:60495"
	destAddress   = "destAddress"
	proxyAddress  = "proxyAddress"
)

func init() {
	const pac = `function matched(url, host) {
    return %s;
}

function FindProxyForURL(url, host) {
    if (matched(url, host)) {
        return "PROXY %s; SOCKS %s; DIRECT";
    }
    return "DIRECT";
}`
	http.HandleFunc("/pac", func(w http.ResponseWriter, r *http.Request) {
		address := menuet.Defaults().String(destAddress)
		matched := menuet.Defaults().String(proxyAddress)
		var m []string
		for _, s := range strings.Split(matched, ",") {
			m = append(m, fmt.Sprintf("shExpMatch(host, %q)", s))
		}
		content := fmt.Sprintf(pac, strings.Join(m, " || "), address, address)
		println(content)
		_, _ = w.Write([]byte(content))
	})
}

func startProxy() bool {
	if server != nil {
		return true
	}

	server = &http.Server{Addr: listenAddress}
	go server.ListenAndServe()
	notification("OVpn", "OVpn is started", "you can click proxy to copy address, then setting with system.")
	return true
}

func proxyItem(items []menuet.MenuItem) menuet.MenuItem {
	text := "Testing proxy - "
	placeholder := "eg: 127.0.0.1:8080"
	proxy := menuet.MenuItem{}

	// default click handler is closing event
	clicked := func() {
		if menuet.Defaults().String(destAddress) == "" {
			clicked := menuet.App().Alert(menuet.Alert{
				MessageText:     "Please enter the proxy dest address",
				InformativeText: "No proxy address can not start",
				Buttons:         []string{"OK", "Cancel"},
				Inputs:          []string{placeholder},
			})

			if clicked.Button == 1 || clicked.Inputs[0] == "" {
				return
			}

			menuet.Defaults().SetString(destAddress, clicked.Inputs[0])
		}

		if server != nil {
			_ = server.Shutdown(context.Background())
			server = nil
		}
		opening = !opening
		if len(items) < 1 {
			return
		}
		items[0] = proxy
	}

	proxy.Text = text + "Stopped"
	proxy.Clicked = clicked

	if opening {
		if !startProxy() {
			return proxy
		}
		proxy.Text = text + "Running"
		proxy.FontWeight = menuet.WeightBold
		proxy.Clicked = func() {
			err := clipboard.WriteAll("http://" + listenAddress + "/pac")
			if err != nil {
				notification(
					"Copy Testing proxy address error",
					"Subtitle",
					err.Error(),
				)
			}
		}
		proxy.Children = func() []menuet.MenuItem {
			da := menuet.Defaults().String(destAddress)
			if da != "" {
				placeholder = da
			}
			subItems := []menuet.MenuItem{
				{
					Text:       da,
					FontSize:   12,
					FontWeight: menuet.WeightBlack,
					Clicked: func() {
						clicked := menuet.App().Alert(menuet.Alert{
							MessageText:     "Wants to update proxy address?",
							InformativeText: "Please re-enter a valid address in the input box",
							Buttons:         []string{"OK", "Cancel"},
							Inputs:          []string{placeholder},
						})

						if clicked.Button == 1 || clicked.Inputs[0] == "" {
							return
						}

						menuet.Defaults().SetString(destAddress, clicked.Inputs[0])
					},
				},
				{Type: menuet.Separator},
			}

			address := menuet.Defaults().String(proxyAddress)
			if address != "" {
				sources := strings.Split(address, ",")
				for i, s := range sources {
					idx := i
					name := s
					subItems = append(subItems, menuet.MenuItem{
						Text:     name,
						FontSize: 12,
						Children: func() []menuet.MenuItem {
							return []menuet.MenuItem{
								{Text: "Delete", Clicked: func() {
									choose := menuet.App().Alert(menuet.Alert{
										MessageText:     "Do you want to delete this address?",
										InformativeText: name,
										Buttons:         []string{"OK", "Cancel"},
									})
									if choose.Button == 0 {
										length := len(sources)
										var newer []string
										if idx == length-1 {
											newer = sources[:length-1]
										} else if idx < length-1 {
											newer = append(sources[:idx], sources[idx+1:]...)
										}
										menuet.Defaults().SetString(proxyAddress, strings.Join(newer, ","))
										fmt.Println(strings.Join(newer, ","))
									}
								}},
							}
						},
					})
				}
			}

			subItems = append(
				subItems,
				menuet.MenuItem{Text: "Add regex", Clicked: func() {
					result := menuet.App().Alert(menuet.Alert{
						MessageText:     "Add a new regex",
						InformativeText: "this regex will be proxy",
						Buttons:         []string{"OK", "Cancel"},
						Inputs:          []string{"Please input valid regex"},
					})

					if result.Button == 1 || result.Inputs[0] == "" {
						return
					}

					content := result.Inputs[0]
					if address != "" {
						content = strings.Join([]string{address, content}, ",")
					}

					menuet.Defaults().SetString(proxyAddress, content)
				}},
				menuet.MenuItem{Type: menuet.Separator},
				menuet.MenuItem{Text: "Stop", Clicked: clicked},
			)
			return subItems
		}
	}
	return proxy
}
