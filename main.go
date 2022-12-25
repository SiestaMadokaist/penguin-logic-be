package main

import (
	"image/png"
	"os"
	"project-penguin-logic/modules/DepotRecognition"
)

func main() {
	imgFile, err := os.Open("./testdata/input/1.png")
	if (err != nil) { panic(err); }
	img, err := png.Decode(imgFile)
	if (err != nil) { panic(err); }
	depotImage := DepotRecognition.NewDepot(img)
	depotEdge := depotImage.GetEdge()
	ioWriter, err := os.Create("./testdata/output/1.edge.png")
	if (err != nil) { panic(err); }
	png.Encode(ioWriter, depotEdge.Image)
}
