package DepotRecognition

import (
	"image"
	"image/color"
	"math"
)

func sumGray(h color.RGBA, v color.RGBA) int {
	// r := float64(p.R) + float64(q.R);
	// g := float64(p.G) + float64(q.G);
	// b := float64(p.B) + float64(q.B);
	r := math.Sqrt(float64((h.R * h.R) + (v.R * v.R)))
	g := math.Sqrt(float64((h.G * h.G) + (v.G * v.G)))
	b := math.Sqrt(float64((h.B * h.B) + (v.B * v.B)))
	gray := math.Sqrt((r * r) + (g * g) + (b * b))
	return int(gray)
}
func (depot *DepotImage) GetEdge() (DepotImage) {
	SobelH := [3][3]int {
		{ 1, 2, 1 },
		{ 0, 0, 0 },
		{ -1, -2, -1 },
	}

	SobelV := [3][3]int {
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}

	edgeH := convolute(depot.Image, SobelH)
	edgeV := convolute(depot.Image, SobelV)

	bounds := depot.Image.Bounds();
	edge := image.NewGray(image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y))
	minX := edge.Bounds().Min.X
	minY := edge.Bounds().Min.Y;

	for xOff := 0; xOff < edge.Bounds().Dx(); xOff++ {
		for yOff := 0; yOff < edge.Bounds().Dy(); yOff++ {
		  x	:= xOff + minX
			y := yOff + minY
			rgbaV := edgeV.RGBAAt(x, y)
			rgbaH := edgeH.RGBAAt(x, y)
			grayLevel := sumGray(rgbaV, rgbaH)
			// fmt.Printf("x: %d y: %d start: (%d, %d) H: %v V: %v G: %v\n", x, y, minX, minY , rgbaH, rgbaV, grayLevel);
			var whiteLevel int;
			if (grayLevel > 13) { whiteLevel = 255 } else { whiteLevel = 0 }
			edge.SetGray(x, y, color.Gray{ Y: uint8(whiteLevel) })
		}
	}
	return DepotImage{Image: edge, Progress: EdgeDetection}
}

func convolute(img image.Image, mask [3][3]int) *image.RGBA {
	bound := img.Bounds();
	minX := bound.Min.X
	minY := bound.Min.Y;
	edgeDetected := image.NewRGBA(image.Rect(bound.Min.X, bound.Min.Y, bound.Max.X, bound.Max.Y))
	for xOff := 1; xOff < img.Bounds().Dx() - 1; xOff++ {
		for yOff := 1; yOff < img.Bounds().Dy() -1; yOff++ {
			x := minX + xOff
			y := minY + yOff;
			r1 := float64(0)
			g1 := float64(0)
			b1 := float64(0)
			for yMask := -1; yMask < 2; yMask++ {
				for xMask := -1; xMask < 2; xMask++ {
					x1 := x + xMask
					y1 := y + yMask
					r, g, b, _ := img.At(x1, y1).RGBA()
					multiplier := float64(mask[yMask + 1][xMask + 1])
					r0 := float64(r / 256)
					g0 := float64(g / 256)
					b0 := float64(b / 256)
					r1 = r1 + (r0 * multiplier)
					g1 = g1 + (g0 * multiplier)
					b1 = b1 + (b0 * multiplier)
				}
			}
			r2 := uint8(math.Max(math.Ceil(r1 / 9), 0))
			g2 := uint8(math.Max(math.Ceil(g1 / 9), 0))
			b2 := uint8(math.Max(math.Ceil(b1 / 9), 0))
			// fmt.Println(x, y, r2, g2, b2);
			edgeDetected.SetRGBA(x, y, color.RGBA{ R: r2, G: g2, B: b2, A: 255})
		}
	}
	return edgeDetected
}