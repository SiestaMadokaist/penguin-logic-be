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
	// bounds := DepotRecognition.GetBounds()
	// fmt.Printf("%v", bounds)
	depotImage := DepotRecognition.NewDepot(img)
	// crop := depotImage.Crop(320, 20, 160)
	depotImage.GetEdge().Save("test.png")
	// depotImage.GetEdge().Save("test.png")
	// crop := depotImage.Crop(396, 450, 160);
	// crop.GetEdge().Save("test.png")
	// println(r, g, b, a)
	// crop.GetEdge().Save("test.png");
	// depotImage.ItemDetectedAt(426, 470, 10);
	// crop.GetEdge();
	// depotImage.GetEdge().Save("edge.png")
	// for i := 0; i < 1; i++ {
	// 	x := 320
	// 	y := 20
	// 	// x := rand.Int() % (depotImage.Image.Bounds().Dx() - 160)
	// 	// y := rand.Int() % (depotImage.Image.Bounds().Dy() - 160)
	// 	img, conf := depotImage.ItemDetectedAt(x, y, 160)
	// 	// fmt.Printf("\n------\n")
	// 	// if (conf > 0.9) {
	// 		// rejections := img.GetEdge().Rejections()
	// 		// fmt.Printf("Location: %d-%d %d %f %d\n", i, x, y, conf, len(rejections))
	// 		name := fmt.Sprintf("./testdata/output/randomcheck/%d-%d-%d-%f.png",i, x, y, conf)
	// 		// for _, r := range rejections {
	// 		// 	fmt.Printf("%v\n", r);
	// 		// }
	// 		fmt.Printf("------\n")
	// 		img.Save(name)
	// 	// }
	// }
	// print(depotEdge)
	// ioWriter, err := os.Create("./testdata/output/1.edge2.png")
	// if (err != nil) { panic(err); }
	// png.Encode(ioWriter, i)
}
