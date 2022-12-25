package main

import (
	"fmt"
	"image/png"
	"math/rand"
	"os"
	"project-penguin-logic/modules/DepotRecognition"
)

func main() {
	imgFile, err := os.Open("./testdata/input/1.png")
	if (err != nil) { panic(err); }
	img, err := png.Decode(imgFile)
	if (err != nil) { panic(err); }
	// bounds := DepotRecognition.GetBounds()
	// fmt.Printf("%v", bounds)
	depotImage := DepotRecognition.NewDepot(img)
	for i := 0; i < 1000; i++ {
		x := rand.Int() % (depotImage.Image.Bounds().Dx() - 160)
		y := rand.Int() % (depotImage.Image.Bounds().Dy() - 160)
		img, conf := depotImage.ItemDetectedAt(x, y, 160)
		if (conf > 0.9) {
			name := fmt.Sprintf("./testdata/output/randomcheck/%d-%d-%d-%f.png",i, x, y, conf)
			println(i, x, y, conf)
			img.Save(name)
		}
	}
	// print(depotEdge)
	// ioWriter, err := os.Create("./testdata/output/1.edge2.png")
	// if (err != nil) { panic(err); }
	// png.Encode(ioWriter, i)
}
