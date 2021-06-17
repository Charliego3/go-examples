package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
	"github.com/progrium/macdriver/webkit"
	"log"
	"runtime"
)

//go:embed index.html
var defaultIndex []byte

func main() {
	spacesFlag := flag.Bool("spaces", true, "appear on all spaces")
	flag.Parse()

	runtime.LockOSThread()
	//var err error
	//
	//usr, err := user.Current()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//dir := filepath.Join(usr.HomeDir, ".topframe")
	//os.MkdirAll(dir, 0755)
	//
	//if _, err := os.Stat(filepath.Join(dir, "index.html")); os.IsNotExist(err) {
	//	ioutil.WriteFile(filepath.Join(dir, "index.html"), defaultIndex, 0644)
	//}
	//
	//srv := http.Server{
	//	Handler: http.FileServer(http.Dir(dir)),
	//}
	//
	//ln, err := net.Listen("tcp", ":0")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//go srv.Serve(ln)

	//fw := watcher.New()
	//if err := fw.AddRecursive(dir); err != nil {
	//	log.Fatal(err)
	//}

	//go fw.Start(400 * time.Millisecond)

	app := cocoa.NSApp_WithDidLaunch(func(notification objc.Object) {
		config := webkit.WKWebViewConfiguration_New()
		config.Preferences().SetValueForKey(core.False, core.String("developerExtrasEnabled"))

		wv := webkit.WKWebView_Init(cocoa.NSScreen_Main().Frame(), config)
		wv.SetOpaque(false)
		wv.SetBackgroundColor(cocoa.NSColor_Clear())
		wv.SetValueForKey(core.False, core.String("drawsBackground"))

		//url := core.URL(fmt.Sprintf("http://localhost:%d", ln.Addr().(*net.TCPAddr).Port))
		url := core.URL("https://www.zb.today/cn/kline/dot_qc")
		req := core.NSURLRequest_Init(url)
		wv.LoadRequest(req)

		frame := cocoa.NSScreen_Main().Frame()
		fmt.Println(frame.String())
		//w := cocoa.NSWindow_Init(cocoa.NSScreen_Main().Frame(),
		w := cocoa.NSWindow_Init(core.NSRect{
			Origin: core.NSPoint{
				X: frame.Size.Width - 900,
				Y: frame.Size.Height - 395,
			},
			Size: core.NSSize{
				Width:  float64(900),
				Height: float64(370),
			},
		},
			cocoa.NSResizableWindowMask|cocoa.NSTitledWindowMask|cocoa.NSClosableWindowMask,
			//cocoa.NSTexturedBackgroundWindowMask|cocoa.NSMiniaturizableWindowMask,
			cocoa.NSBackingStoreBuffered, false)
		w.SetContentView(wv)
		w.SetBackgroundColor(cocoa.NSColor_Clear())
		w.SetOpaque(true)
		w.SetTitleVisibility(cocoa.NSWindowTitleHidden)
		w.SetTitlebarAppearsTransparent(false)
		w.SetIgnoresMouseEvents(false)
		w.SetMovableByWindowBackground(true)
		w.SetLevel(cocoa.NSMainMenuWindowLevel + 2)
		w.MakeKeyAndOrderFront(w)
		if *spacesFlag {
			w.SetCollectionBehavior(cocoa.NSWindowCollectionBehaviorCanJoinAllSpaces)
			//w.SetCollectionBehavior(cocoa.NSWindowCollectionBehaviorIgnoresCycle)
		}

		//events := make(chan cocoa.NSEvent)
		//go func() {
		//	for e := range events {
		//		if c, _ := e.KeyCode(); c == 100 {
		//			if w.IgnoresMouseEvents() {
		//				fmt.Println("Mouse events on")
		//				w.SetIgnoresMouseEvents(false)
		//			} else {
		//				fmt.Println("Mouse events off")
		//				w.SetIgnoresMouseEvents(true)
		//			}
		//		}
		//		e.Release()
		//	}
		//}()
		//cocoa.NSEvent_GlobalMonitorMatchingMask(cocoa.NSEventMaskKeyDown, events)

		//go func() {
		//	for {
		//		select {
		//		case event := <-fw.Event:
		//			if event.IsDir() {
		//				continue
		//			}
		//			wv.Reload(nil)
		//		case <-fw.Closed:
		//			return
		//		}
		//	}
		//}()
	})

	app.SetActivationPolicy(cocoa.NSApplicationActivationPolicyRegular)
	//app.ActivateIgnoringOtherApps(true)

	log.Printf("topframe 0.1.0 by progrium\n")
	app.Run()
}
