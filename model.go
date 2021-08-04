package main

import (
	"errors"
)

// ---
type Cmd struct {
	Uuid  string
	Name  string
	Cdate int64
}

// ---
type AddStrokeContents struct {
	Uuid      string //strokeObjectUuid
	Color     int64
	GroupUuid string
	Pts       []float64
}

type AddStrokeCmd struct {
	Contents AddStrokeContents
}

// ---
type DeleteStrokesContents struct {
	Uuids []string
}

type DeleteStrokesCmd struct {
	Contents DeleteStrokesContents
}

// ---
type MoveStrokesContents struct {
	Uuids []string
	Pts   []float64 //transform information
}

type MoveStrokesCmd struct {
	Contents MoveStrokesContents
}

// ---
type ResizeGroupContents struct {
	ScaleX      float64
	ScaleY      float64
	StrokeUuids []string
}

type ResizeGroupCmd struct {
	Contents ResizeGroupContents
}

// ---
type StrokeObject struct {
	PageUuid           string
	Uuid               string
	Pts                []float64
	Color              int64
	GroupUuid          string
	LogicalStrokeWidth int64
}

type Db struct {
	StrokeObjectSlice []StrokeObject
}

func (db *Db) Add(strokeObject StrokeObject) {
	db.StrokeObjectSlice = append(db.StrokeObjectSlice, strokeObject)
}

func (db *Db) Remove(strokeObjectUuid string) {
	result := []StrokeObject{}
	for i := range db.StrokeObjectSlice {
		strokeObject := db.StrokeObjectSlice[i]
		if strokeObject.Uuid != strokeObjectUuid {
			result = append(result, strokeObject)
		}
	}
	db.StrokeObjectSlice = result
}

func (db *Db) Get(strokeObjectUuid string) (StrokeObject, error) {
	var retVal StrokeObject
	found := false
	for i := range db.StrokeObjectSlice {
		strokeObject := db.StrokeObjectSlice[i]
		if strokeObject.Uuid == strokeObjectUuid {
			retVal = strokeObject
			found = true
		}
	}

	if found == true {
		return retVal, nil
	} else {
		return *new(StrokeObject), errors.New("Not Found")
	}
}

func toRectangle(strokeObjectSlice []StrokeObject) Rectangle {
	maxLen := len(strokeObjectSlice)
	if maxLen > 0 {
		xlist := []float64{}
		ylist := []float64{}

		for i := range strokeObjectSlice {
			strokeObject := strokeObjectSlice[i]
			pts := strokeObject.Pts
			maxLenHalf := len(pts) / 2
			for j := 0; j < maxLenHalf; j++ {
				xindex := j * 2
				yindex := xindex + 1
				x := pts[xindex]
				y := pts[yindex]
				xlist = append(xlist, x)
				ylist = append(ylist, y)
			}
		}

		if len(xlist) > 0 && len(ylist) > 0 {
			// find min x as left and max x as right:
			left := xlist[0]
			right := xlist[0]
			for i := range xlist {
				x := xlist[i]
				if x < left {
					left = x
				}
				if right < x {
					right = x
				}
			}

			// find min y as top and max y as bottom:
			top := ylist[0]
			bottom := ylist[0]
			for i := range ylist {
				y := ylist[i]
				if y < top {
					top = y
				}
				if bottom < y {
					bottom = y
				}
			}
			return Rectangle{left, top, right, bottom}
		} else {
			return Rectangle{0.0, 0.0, 0.0, 0.0}
		}

	} else {
		return Rectangle{0.0, 0.0, 0.0, 0.0}
	}
}

func (db *Db) toRectangle() Rectangle {
	return toRectangle(db.StrokeObjectSlice)
}
