package main

import (
	"bytes"
	"fmt"
	"image/png"
	"log"
	"regexp"
	"time"

	"github.com/google/gousb"
	"github.com/kbinani/screenshot"
	"github.com/otiai10/gosseract"
)

var endpointDevice *gousb.OutEndpoint

func main() {
	endpointDevice = getDualSenseEndpoint(gousb.ID(0x054c), gousb.ID(0x0ce6))
	watchScreen()
}

func getDualSenseEndpoint(vendor gousb.ID, product gousb.ID) *gousb.OutEndpoint {
	ctx := gousb.NewContext()
	ctx.Debug(4)
	controller, err := ctx.OpenDeviceWithVIDPID(vendor, product)
	if err != nil {
		log.Fatalf("Nao foi possivel abrir dispositivo: %v", err)
	}
	defer controller.Close()

	cfg, err := controller.Config(1)
	if err != nil {
		log.Fatalf("%s.Config(1): %v", controller, err)
	}
	defer cfg.Close()

	intf, err := cfg.Interface(0, 0)
	if err != nil {
		log.Fatalf("%s.Interface(0, 0): %v", cfg, err)
	}
	defer intf.Close()

	epOut, err := intf.OutEndpoint(7)
	if err != nil {
		log.Fatalf("%s.OutEndpoint(7): %v", intf, err)
	}

	return epOut
}

func pressButton(button string) {
	switch button {
	case "cross":
		{
			state := make([]byte, 79)
			state[0] = 0x7E
			endpointDevice.Write(state)
		}
	case "square":
		{
			state := make([]byte, 79)
			state[0] = 0x05
			state[1] = 0x07
			endpointDevice.Write(state)
		}
	case "triangle":
		{
			state := make([]byte, 79)
			state[0] = 0x05
			state[1] = 0x07
			endpointDevice.Write(state)
		}
	case "circle":
		{
			state := make([]byte, 79)
			state[0] = 0x05
			state[1] = 0x07
			endpointDevice.Write(state)
		}
	case "options":
		{
			state := make([]byte, 79)
			state[0] = 0x05
			state[1] = 0x07
			endpointDevice.Write(state)
		}
	}
}

func watchScreen() {
	var validText = regexp.MustCompile(`(?i)\bOnline\b`)

	client := gosseract.NewClient()
	defer client.Close()

	for {
		bounds := screenshot.GetDisplayBounds(0)
		img, _ := screenshot.CaptureRect(bounds)
		buf := new(bytes.Buffer)
		png.Encode(buf, img)

		client.SetLanguage("por", "eng")
		client.SetImageFromBytes(buf.Bytes())
		text, _ := client.Text()

		// pressButton("cross")

		if validText.MatchString(text) {
			fmt.Println(text)
			pressButton("cross")
		}

		time.Sleep(1 * time.Second)
	}
}
