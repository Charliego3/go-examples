package main

import (
	"fmt"
	usbdrivedetector "github.com/deepakjois/gousbdrivedetector"
	"github.com/google/gousb"
	gep "github.com/mect/go-escpos"
	"github.com/pzl/usb"
	"github.com/sirupsen/logrus"
	"github.com/tarm/serial"
	"log"
	"time"
)

func main() {
	/*f, err := os.OpenFile("/dev/ttys000", os.O_RDWR, 0755)
	if err != nil {
		logrus.Fatal(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	p := escpos.New(w)

	p.Verbose = true

	p.Init()
	p.Beep(4)
	p.SetFontSize(2, 3)
	p.SetFont("B")
	p.SetReverse(0)
	p.WriteGBK("简体字转繁体字")
	p.SetFont("C")
	p.Write("test2")

	p.SetEmphasize(1)
	p.Write("hello")
	p.Formfeed()
	// png, _ = qrcode.Encode("https://www.bing.com", qrcode.Low, 256)
	// img, _, _ := image.Decode(bytes.NewReader(png))
	// p.SetAlign("center")
	// p.PrintImage(img)
	p.SetUnderline(1)
	p.SetFontSize(4, 4)
	p.Write("hello")
	p.SetReverse(1)
	p.SetFontSize(2, 4)
	p.Write("hello")
	p.FormfeedN(10)
	p.SetAlign("center")
	p.Write("test")
	p.Linefeed()

	p.SetEmphasize(0)
	p.SetReverse(0)
	p.SetFontSize(2, 2)
	p.SetUnderline(0)
	data := [][]string{
		[]string{"充值", "The Good", "500"},
		[]string{"找零", "The Ruby", "288"},
		[]string{"应收", "The Ugly", "120"},
		[]string{"实收", "The Gopher", "800"},
	}
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Name", "Sign", "Rating"})

	for _, v := range data {
		table.Append(v)
	}
	table.SetAutoFormatHeaders(true)
	table.SetAutoMergeCells(true)
	table.SetAutoWrapText(true)
	table.SetBorder(true)
	table.Render() // Send output
	p.WriteGBK(tableString.String())
	p.Linefeed()
	p.Write("test")
	p.FormfeedD(1)

	p.Cut()

	w.Flush()*/

	// goUSB()

	if drives, err := usbdrivedetector.Detect(); err == nil {
		fmt.Printf("%d USB Devices Found\n", len(drives))
		for _, d := range drives {
			fmt.Println(d)
		}
	} else {
		fmt.Println(err)
	}
}

func pzlusb() {
	devices, err := usb.List()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("Devices: %#v", devices)
}

func libusbtest() {
	// ctx, err := libusb.NewContext()
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	//
	// ctx.SetDebug(libusb.LogLevelDebug)
	// devices, err := ctx.GetDeviceList()
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	//
	// logrus.Infof("Devices: %#v", devices)
}

func goUSB() {
	ctx := gousb.NewContext()
	defer ctx.Close()

	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		logrus.Infof("%#vs", desc)
		for _, v := range desc.Configs {
			for _, i := range v.Interfaces {
				for _, s := range i.AltSettings {
					// logrus.Infof("Settings: %#v", s)
					if s.Class == gousb.ClassPrinter {
						return true
					}
				}
			}
		}

		return false
	})

	if err != nil {
		logrus.Fatal(err)
	}

	for _, dev := range devs {
		// logrus.Info(dev.Desc.Class)

		num, err := dev.ActiveConfigNum()
		if err != nil {
			return
		}

		product, err := dev.Product()
		if err != nil {
			logrus.Error(err)
			continue
		}

		logrus.Infof("Product: %s", product)

		for i := 0; i < num; i++ {
			number, err := dev.SerialNumber()
			if err != nil {
				logrus.Error(err)
				continue
			}

			logrus.Infof("SerialNumber: %s", number)
		}
	}

	dev := devs[0]
	logrus.Infof("Vendor: %s, Product: %s", dev.Desc.Vendor.String(), dev.Desc.Product)
	device, err := ctx.OpenDeviceWithVIDPID(1137, 85)
	// device, err := ctx.OpenDeviceWithVIDPID(dev.Desc.Vendor, dev.Desc.Product)
	if err != nil {
		logrus.Fatal(err)
	}
	defer device.Close()

	product, err := device.Product()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("Device Product: %v", product)

	// intf, done, err := device.DefaultInterface()
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	// defer done()

	// var oe *gousb.OutEndpoint
	// cfg, err := device.Config(1)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// intf, err := cfg.Interface(0, 0)
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	//
	// defer intf.Close()
	// // for _, desc := range dev.Desc.Configs[0].Interfaces[0].AltSettings[0].Endpoints {
	// // 	if desc.Direction == gousb.EndpointDirectionOut {
	// oe, _ = intf.OutEndpoint(2)
	// // 		break
	// // 	}
	// // }
	//
	// w := bufio.NewWriter(oe)
	// p := escpos.New(w)
	//
	// // p.Verbose = true
	// p.SetFont("B")
	//
	// p.Init()
	// // p.Beep(4)
	// p.SetFontSize(2, 3)
	// p.SetReverse(0)
	// p.WriteGBK("简体字转繁体字")
	// p.SetFontSize(1, 1)
	// p.Formfeed()
	// p.Write("test2")
	//
	// p.SetEmphasize(1)
	// p.Write("hello")
	// p.Formfeed()
	// // png, _ = qrcode.Encode("https://www.bing.com", qrcode.Low, 256)
	// // img, _, _ := image.Decode(bytes.NewReader(png))
	// // p.SetAlign("center")
	// // p.PrintImage(img)
	//
	// p.SetAlign("left")
	// p.SetUnderline(1)
	// p.Write("hello")
	// p.SetMoveX(20)
	// p.Write("hello")
	// p.SetAlign("right")
	// p.Write("right")
	// p.Linefeed()
	// p.SetAlign("left")
	//
	// p.SetEmphasize(0)
	// p.SetReverse(0)
	// p.SetUnderline(0)
	//
	// data := [][]string{
	// 	[]string{"简体字转繁体字简体", "500"},
	// 	[]string{"The Ruby", "288"},
	// 	[]string{"The Ugly", "120"},
	// 	[]string{"The Gopher", "800"},
	// }
	// tableString := &strings.Builder{}
	// table := tablewriter.NewWriter(tableString)
	// // table.SetColMinWidth(0, 15)
	// table.SetHeader([]string{"名称", "数量"})
	//
	// for _, v := range data {
	// 	table.Append(v)
	// }
	// // table.SetCaption(true, "这是标题")
	// // table.SetRowSeparator("-")
	// // table.SetAlignment(tablewriter.ALIGN_LEFT)
	// table.SetAutoFormatHeaders(true)
	// table.SetAutoMergeCells(true)
	// table.SetAutoWrapText(true)
	// table.SetBorder(true)
	// table.Render() // Send output
	// p.SetFontSize(1, 1)
	// p.WriteGBK(tableString.String())
	// p.Linefeed()
	// p.Write("test")
	// p.FormfeedD(1)
	//
	// p.Cut()
	//
	// w.Flush()
}

func goescpos() {
	p, err := gep.NewUSBPrinterByPath("") // empry string will do a self discovery
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("Printer: %#v", p)
}

func escposgo() {
	// conn, err := conntion.NewNetConntion("192.168.1.71:9100")
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	//
	// ep := escpos.NewEscpos(conn)
	// ep.SetFont("C")
	// ep.Init()
	// ep.WriteGbk("测试打印")
	// ep.Linefeed()
	// ep.Write("Regular")
	// ep.FormfeedN(1)
	// ep.Cut()
	// ep.End()
}

func goserial() {
	c := &serial.Config{Name: "tty", Baud: 19200, ReadTimeout: time.Second * 5}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	n, err := s.Write([]byte("test"))
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 128)
	n, err = s.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%q", buf[:n])
}
