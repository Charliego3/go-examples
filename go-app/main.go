package main

import "github.com/webview/webview"

func main() {
	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle("JupyterLab")
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("http://localhost:8888/?token=5d5ef3ebde4308cb4c4a126997bef257b4d9f96be03a0cb4")
	w.Run()

	//h := &app.Handler{
	//	Title:  "Hello Demo",
	//	Author: "Maxence Charriere",
	//}
	//
	//if err := http.ListenAndServe(":8888", h); err != nil {
	//	panic(err)
	//}
	//
	//
	//// Initialize astilectron
	//var a, _ = astilectron.New(nil, astilectron.Options{
	//	AppName: "Jupyter",
	//	//AppIconDefaultPath: "<your .png icon>", // If path is relative, it must be relative to the data directory
	//	//AppIconDarwinPath:  "<your .icns icon>", // Same here
	//	//BaseDirectoryPath: "<where you want the provisioner to install the dependencies>",
	//	VersionAstilectron: "0.33.0",
	//	VersionElectron: "4.0.1",
	//})
	//defer a.Close()
	//
	//a.HandleSignals()
	//
	//// Create a new window
	//var w, _ = a.NewWindow("http://localhost:8888/?token=5d5ef3ebde4308cb4c4a126997bef257b4d9f96be03a0cb4", &astilectron.WindowOptions{
	//	Center: astikit.BoolPtr(true),
	//	Height: astikit.IntPtr(600),
	//	Width:  astikit.IntPtr(600),
	//})
	////w.Show()
	//w.Create()
	//
	//// Start astilectron
	//a.Start()
	//
	//// Blocking pattern
	//a.Wait()
}
