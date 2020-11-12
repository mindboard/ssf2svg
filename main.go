package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func createTransformMatrix(scaleX float64, scaleY float64, oldGroupBounds Rectangle) matrix2d {
	oldGroupCenterX := oldGroupBounds.left + (oldGroupBounds.right-oldGroupBounds.left)*0.5
	oldGroupCenterY := oldGroupBounds.top + (oldGroupBounds.bottom-oldGroupBounds.top)*0.5

	m0 := matrix2d{
		1, 0, oldGroupCenterX * (-1),
		0, 1, oldGroupCenterY * (-1),
		0, 0, 1}

	m1 := matrix2d{
		scaleX, 0, 0,
		0, scaleY, 0,
		0, 0, 1}

	m2 := matrix2d{
		1, 0, oldGroupCenterX,
		0, 1, oldGroupCenterY,
		0, 0, 1}

	return m2.multiply(m1.multiply(m0))
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func main() {
	var styleFilePath string
	if len(os.Args) > 1 {
		styleFilePath = os.Args[1]
	}

	var style *Style
	if Exists(styleFilePath) {
		style = createStyle(styleFilePath)
	} else {
		style = createDefaultStyle()
	}

	bs, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	var buf *bytes.Buffer = bytes.NewBuffer(bs)

	zr, err := gzip.NewReader(buf)
	if err != nil {
		log.Fatal(err)
	}

	var db Db

	scanner := bufio.NewScanner(zr)
	for scanner.Scan() {
		line := scanner.Text()

		cmd0 := new(Cmd)
		err := json.Unmarshal([]byte(line), cmd0)
		if err != nil {
			log.Fatal(err)
		}

		if cmd0.Name == "ADD_STROKE" {
			cmd1 := new(AddStrokeCmd)
			err := json.Unmarshal([]byte(line), cmd1)
			if err != nil {
				log.Fatal(err)
			}

			strokeObject := new(StrokeObject)
			strokeObject.PageUuid = "0" // default value
			strokeObject.Uuid = cmd1.Contents.Uuid
			strokeObject.Pts = cmd1.Contents.Pts
			strokeObject.Color = cmd1.Contents.Color
			strokeObject.GroupUuid = cmd1.Contents.GroupUuid
			strokeObject.LogicalStrokeWidth = 0 // default value

			db.Add(*strokeObject)
		}

		if cmd0.Name == "DELETE_STROKES" {
			cmd1 := new(DeleteStrokesCmd)
			err := json.Unmarshal([]byte(line), cmd1)
			if err != nil {
				log.Fatal(err)
			}

			for i := range cmd1.Contents.Uuids {
				strokeObjectUuid := cmd1.Contents.Uuids[i]
				db.Remove(strokeObjectUuid)
			}
		}

		if cmd0.Name == "MOVE_STROKES" {
			cmd1 := new(MoveStrokesCmd)
			err := json.Unmarshal([]byte(line), cmd1)
			if err != nil {
				log.Fatal(err)
			}

			var startX float64
			var startY float64
			var stopX float64
			var stopY float64

			if len(cmd1.Contents.Pts) > 3 {
				startX = cmd1.Contents.Pts[0]
				startY = cmd1.Contents.Pts[1]
				stopX = cmd1.Contents.Pts[2]
				stopY = cmd1.Contents.Pts[3]

				translateMatrix := matrix2d{
					1, 0, (stopX - startX),
					0, 1, (stopY - startY),
					0, 0, 1}

				for i := range cmd1.Contents.Uuids {
					strokeObjectUuid := cmd1.Contents.Uuids[i]
					strokeObject, err := db.Get(strokeObjectUuid)
					if err == nil {
						// 1)
						db.Remove(strokeObjectUuid)

						// 2)
						myStrokeObject := new(StrokeObject)
						myStrokeObject.PageUuid = strokeObject.PageUuid
						myStrokeObject.Uuid = strokeObjectUuid
						myStrokeObject.Pts = mapPoints(translateMatrix, strokeObject.Pts)
						myStrokeObject.Color = strokeObject.Color
						myStrokeObject.GroupUuid = strokeObject.GroupUuid
						myStrokeObject.LogicalStrokeWidth = strokeObject.LogicalStrokeWidth

						db.Add(*myStrokeObject)
					}
				}
			}
		}

		if cmd0.Name == "RESIZE_GROUP" {
			cmd1 := new(ResizeGroupCmd)
			err := json.Unmarshal([]byte(line), cmd1)
			if err != nil {
				log.Fatal(err)
			}

			// 1)
			scaleX := cmd1.Contents.ScaleX
			scaleY := cmd1.Contents.ScaleY

			strokeObjectSlice := []StrokeObject{}
			for i := range cmd1.Contents.StrokeUuids {
				strokeObjectUuid := cmd1.Contents.StrokeUuids[i]
				strokeObject, err := db.Get(strokeObjectUuid)
				if err == nil {
					strokeObjectSlice = append(strokeObjectSlice, strokeObject)
				}
			}

			groupRectangle := toRectangle(strokeObjectSlice)

			transformMatrixValues := createTransformMatrix(scaleX, scaleY, groupRectangle)

			// 2)
			for i := range cmd1.Contents.StrokeUuids {
				strokeObjectUuid := cmd1.Contents.StrokeUuids[i]
				strokeObject, err := db.Get(strokeObjectUuid)
				if err == nil {
					// 2-1)
					db.Remove(strokeObjectUuid)

					// 2-2)
					myStrokeObject := new(StrokeObject)
					myStrokeObject.PageUuid = strokeObject.PageUuid
					myStrokeObject.Uuid = strokeObjectUuid
					myStrokeObject.Pts = mapPoints(transformMatrixValues, strokeObject.Pts)
					myStrokeObject.Color = strokeObject.Color
					myStrokeObject.GroupUuid = strokeObject.GroupUuid
					myStrokeObject.LogicalStrokeWidth = strokeObject.LogicalStrokeWidth

					db.Add(*myStrokeObject)
				}
			}
		}
	}

	if err := zr.Close(); err != nil {
		log.Fatal(err)
	}

	svg := createSvg(db, style)

	fmt.Fprintln(os.Stdout, svg)
}
