package main

type point2d struct {
	x float64
	y float64
}

type matrix2d struct {
	// 1st row
	e00 float64
	e01 float64
	e02 float64

	// 2nd row
	e10 float64
	e11 float64
	e12 float64

	// 3rd row
	e20 float64
	e21 float64
	e22 float64
}

func (m0 *matrix2d) multiply(m1 matrix2d) matrix2d {
	a0 := m0.e00
	b0 := m0.e01
	c0 := m0.e02
	d0 := m0.e10
	e0 := m0.e11
	f0 := m0.e12
	g0 := m0.e20
	h0 := m0.e21
	i0 := m0.e22

	a1 := m1.e00
	b1 := m1.e01
	c1 := m1.e02
	d1 := m1.e10
	e1 := m1.e11
	f1 := m1.e12
	g1 := m1.e20
	h1 := m1.e21
	i1 := m1.e22

	// 1st row:
	e00 := a0*a1 + b0*d1 + c0*g1
	e01 := a0*b1 + b0*e1 + c0*h1
	e02 := a0*c1 + b0*f1 + c0*i1

	// 2nd row:
	e10 := d0*a1 + e0*d1 + f0*g1
	e11 := d0*b1 + e0*e1 + f0*h1
	e12 := d0*c1 + e0*f1 + f0*i1

	// 3rd row:
	e20 := g0*a1 + h0*d1 + i0*g1
	e21 := g0*b1 + h0*e1 + i0*h1
	e22 := g0*c1 + h0*f1 + i0*i1

	return matrix2d{
		e00, e01, e02,
		e10, e11, e12,
		e20, e21, e22}
}

func (m *matrix2d) transform(pt point2d) point2d {
	a := m.e00
	b := m.e01
	c := m.e02

	d := m.e10
	e := m.e11
	f := m.e12

	x := pt.x
	y := pt.y
	z := 1.0

	x0 := a*x + b*y + c*z
	y0 := d*x + e*y + f*z

	return point2d{x0, y0}
}

func mapPoints(matrix matrix2d, pts []float64) []float64 {
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

	for i := range xlist {
		x := xlist[i]
		y := ylist[i]

		newPt := matrix.transform(point2d{x, y})
		xlist[i] = newPt.x
		ylist[i] = newPt.y
	}

	var newPts []float64
	for i := range xlist {
		x := xlist[i]
		y := ylist[i]
		newPts = append(newPts, x)
		newPts = append(newPts, y)
	}

	return newPts
}
