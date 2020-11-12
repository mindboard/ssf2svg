package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Color struct {
	R uint8
	G uint8
	B uint8
}

type Style struct {
	Border          bool
	Padding         bool
	FillBackground  bool
	BorderColor     Color
	BorderWidth     float64
	StrokeWidth     float64
	StrokeColorList []Color
	BackgroundColor Color
	CanvasWidth     uint64
}

func createDefaultStyle() *Style {
	style := new(Style)
	style.Border = false
	style.Padding = true
	style.FillBackground = true
	style.BorderColor = Color{R: 27, G: 95, B: 133}
	style.BorderWidth = 1.2
	style.StrokeWidth = 3.2
	style.StrokeColorList = []Color{
		Color{R: 76, G: 96, B: 130},
		Color{R: 232, G: 83, B: 33},
		Color{R: 240, G: 149, B: 31},
		Color{R: 0, G: 159, B: 129},
		Color{R: 177, G: 150, B: 147},
		Color{R: 8, G: 8, B: 8}}
	style.BackgroundColor = Color{R: 253, G: 246, B: 227}
	style.CanvasWidth = 1024

	return style
}

func createStyle(filepath string) *Style {
	f, err := os.Open(filepath)

	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	style := new(Style)

	err2 := json.Unmarshal(bs, style)
	if err2 != nil {
		log.Fatal(err2)
	}

	return style
}
