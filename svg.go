package main

import (
	"bytes"
	"fmt"
)

type Rectangle struct {
	left float64
	top float64
	right float64
	bottom float64
}

func createXYWH(bounds Rectangle) string {
	x := bounds.left
	y := bounds.top
	w := bounds.right - bounds.left
	h := bounds.bottom - bounds.top
	return fmt.Sprintf("x=\"%v\" y=\"%v\" width=\"%v\" height=\"%v\" ", x, y, w, h)
}

func createViewBox(bounds Rectangle) string {
	x := bounds.left
	y := bounds.top
	w := bounds.right - bounds.left
	h := bounds.bottom - bounds.top
	return fmt.Sprintf("viewBox=\"%v %v %v %v\" ", x, y, w, h)
}

func createSvgColor(color Color) string {
	return fmt.Sprintf("rgb(%+v, %+v, %+v)", color.R, color.G, color.B)
}

func createSvgStrokeColor(strokeColor int64, style *Style) string {
	for i := range []int{0, 1, 2, 3, 4, 5} {
		if i < len(style.StrokeColorList) {
			if strokeColor == int64(i) {
				return createSvgColor(style.StrokeColorList[i])
			}
		}
	}

	if len(style.StrokeColorList) > 0 {
		strokeColorObject := style.StrokeColorList[0]
		return fmt.Sprintf("rgb(%+v, %+v, %+v)", strokeColorObject.R, strokeColorObject.G, strokeColorObject.B)
	}

	return "rgb(0,0,0)"
}

func createSvgHeader(bounds Rectangle, style *Style) string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("<?xml version=\"1.0\" encoding=\"utf-8\"?>"))
	buffer.WriteString("\n")
	buffer.WriteString("<!DOCTYPE svg PUBLIC \"-//W3C//DTD SVG 1.1//EN\" \"http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd\">")
	buffer.WriteString("\n")
	buffer.WriteString("<svg version=\"1.1\" xmlns=\"http://www.w3.org/2000/svg\" ")
	buffer.WriteString("\n")
	buffer.WriteString(createXYWH(bounds))
	buffer.WriteString(" ")
	buffer.WriteString(createViewBox(bounds))
	buffer.WriteString(">")
	buffer.WriteString("\n")

	//
	// paint border and background:
	//

	var bgColor string
	if style.FillBackground {
		bgColor = createSvgColor(style.BackgroundColor)
	} else {
		bgColor = "none"
	}

	var borderStrokeColor string
	if style.Border {
		borderStrokeColor = createSvgColor(style.BorderColor)
	} else {
		borderStrokeColor = bgColor
	}

	borderWidth := style.BorderWidth

	buffer.WriteString(fmt.Sprintf("<g style=\"stroke:%s\" stroke-width=\"%v\" fill=\"%s\">", borderStrokeColor, borderWidth, bgColor))
	buffer.WriteString(fmt.Sprintf("<rect %s />", createXYWH(bounds)))
	buffer.WriteString("</g>")
	buffer.WriteString("\n")

	return buffer.String()
}

func createSvgPart(pts []float64, strokeColor int64, style *Style) string {
	var xlist []float64
	var ylist []float64
	for i := range pts {
		if i%2 == 0 {
			// even
			x := pts[i]
			xlist = append(xlist, x)
		}

		if i%2 == 1 {
			// odd
			y := pts[i]
			ylist = append(ylist, y)
		}
	}

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("<g stroke=\"%s\" stroke-width=\"%v\" fill=\"none\">", createSvgStrokeColor(strokeColor, style), style.StrokeWidth))

	lastIndex := len(xlist) - 1

	for i := range xlist {
		if i == 0 {
			x := xlist[i]
			y := ylist[i]

			buffer.WriteString(fmt.Sprintf("<path d=\"M %v %v ", x, y))
		}

		if i > 0 {
			x := xlist[i]
			y := ylist[i]

			if i != lastIndex {
				buffer.WriteString(fmt.Sprintf("L %v %v ", x, y))
			} else {
				buffer.WriteString(fmt.Sprintf("L %v %v\"/>", x, y))
			}
		}
	}

	buffer.WriteString(fmt.Sprintf("</g>\n"))

	return buffer.String()
}

func createSvg(db Db, style *Style) string {
	canvasRectangle := db.toRectangle()
	canvasLeft := canvasRectangle.left
	canvasTop := canvasRectangle.top
	canvasWidth := canvasRectangle.right - canvasRectangle.left
	canvasHeight := canvasRectangle.bottom - canvasRectangle.top

	svgWidth := float64(style.CanvasWidth)
	scale := svgWidth / canvasWidth
	svgHeight := canvasHeight * scale

	var paddingFactor float64 = 1.0
	if style.Padding {
		paddingFactor = 0.9
	}

	matrix0 := matrix2d{
		1, 0, (0.0 - canvasLeft) - canvasWidth/2.0,
		0, 1, (0.0 - canvasTop) - canvasHeight/2.0,
		0, 0, 1}

	matrix1 := matrix2d{
		scale * paddingFactor, 0, 0,
		0, scale * paddingFactor, 0,
		0, 0, 1}

	matrix2 := matrix2d{
		1, 0, (canvasWidth * scale) / 2.0,
		0, 1, (canvasHeight * scale) / 2.0,
		0, 0, 1}

	matrixResult := matrix2.multiply(matrix1.multiply(matrix0))

	bounds := Rectangle{0, 0, svgWidth, svgHeight}

	var buffer bytes.Buffer
	buffer.WriteString(createSvgHeader(bounds, style))

	for i := range db.StrokeObjectSlice {
		strokeObject := db.StrokeObjectSlice[i]

		if len(strokeObject.Pts) > 3 {
			pts := mapPoints(matrixResult, strokeObject.Pts)
			svgCode := createSvgPart(pts, strokeObject.Color, style)
			buffer.WriteString(svgCode)
		}
	}
	buffer.WriteString("</svg>")

	return buffer.String()
}
