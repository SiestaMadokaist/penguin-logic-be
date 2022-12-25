package DepotRecognition

import (
	"image"
	"image/color"
	"math"
)

func sumGray(p color.RGBA, q color.RGBA) int {
	r := math.Sqrt(float64((p.R * p.R) + (q.R * q.R)))
	g := math.Sqrt(float64((p.G * p.G) + (q.G * q.G)))
	b := math.Sqrt(float64((p.B * p.B) + (q.B * q.B)))
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
		{1, 0, -1},
		{2, 0, -2},
		{1, 0, -1},
	}

	edgeH := convolute(depot.Image, SobelH)
	edgeV := convolute(depot.Image, SobelV)

	edge := image.NewGray(image.Rect(0, 0, depot.Image.Bounds().Dx(), depot.Image.Bounds().Dy()))
	
	for x := 0; x < edge.Bounds().Dx(); x++ {
		for y := 0; y < edge.Bounds().Dy(); y++ {
			rgbaV := edgeV.RGBAAt(x, y)
			rgbaH := edgeH.RGBAAt(x, y)
			grayLevel := sumGray(rgbaV, rgbaH)
			var whiteLevel int;
			if (grayLevel > 10) { whiteLevel = 255 } else { whiteLevel = 0}
			edge.SetGray(x, y, color.Gray{ Y: uint8(whiteLevel) })
		}
	}
	return DepotImage{Image: edge, Progress: EdgeDetection}
}

func convolute(img image.Image, mask [3][3]int) *image.RGBA {
	edgeDetected := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))

	for x := 1; x < img.Bounds().Dx() - 1; x++ {
		for y := 1; y < img.Bounds().Dy() -1; y++ {
	// x := 26
	// y := 25
			r1 := int32(0)
			g1 := int32(0)
			b1 := int32(0)
			for yMask := -1; yMask < 2; yMask++ {
				for xMask := -1; xMask < 2; xMask++ {
					x1 := x + xMask
					y1 := y + yMask
					r, g, b, _ := img.At(x1, y1).RGBA()
					multiplier := int32(mask[yMask + 1][xMask + 1])
					r0 := int32(r / 256)
					g0 := int32(g / 256)
					b0 := int32(b / 256)
					r1 = r1 + (r0 * multiplier)
					g1 = g1 + (g0 * multiplier)
					b1 = b1 + (b0 * multiplier)
				}
			}
			r2 := uint8(math.Max(float64(r1 / 9), 0))
			g2 := uint8(math.Max(float64(g1 / 9), 0))
			b2 := uint8(math.Max(float64(b1 / 9), 0))
			edgeDetected.SetRGBA(x, y, color.RGBA{ R: r2, G: g2, B: b2, A: 255})
		}
	}
	return edgeDetected
}