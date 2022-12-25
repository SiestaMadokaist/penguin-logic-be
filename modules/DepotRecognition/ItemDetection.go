package DepotRecognition

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)

type Cropable interface {
	SubImage(r image.Rectangle) image.Image
}

func GetBounds() [][2]int {
	b, err := ioutil.ReadFile("./features/detector.txt")
	if (err != nil) {
		panic(errors.New("cannot open features/detector"))
	}
	boundaries := make([][2]int, 256)
	content := strings.Split(string(b), "\n");
	for i, line := range content {
		v := strings.Split(line, " ")
		low, _ := strconv.ParseInt(v[0], 10, 8)
		high, _ := strconv.ParseInt(v[1], 10, 8)
		boundaries[i] = [2]int{ int(low), int(high) }
	}
	return boundaries;
}

func getConfidence(bounds [][2]int, values []int) float32 {
	inRange := 0
	for i, bound := range bounds {
		lowerbound := bound[0]
		upperbound := bound[1]
		value := values[i];
		if lowerbound <= value && value <= upperbound {
			inRange++
		}
	}
	return float32(inRange) / 256
}

func (depot DepotImage) Crop(x int, y int, size int) DepotImage {
	original, ok := depot.Image.(Cropable)
	if (!ok) { panic("cannot convert image to gray"); }
	cropped := original.SubImage(image.Rect(x, y, x + size, y + size));
	return DepotImage{Image: cropped, Progress: depot.Progress }
}

func (depot DepotImage) getHistogram() []int {
	histogram := make([]int, 256);
	grayImage := depot.Image.(*image.Gray)
	for x := 0; x < depot.Image.Bounds().Dx(); x++ {
		for y := 0; y < depot.Image.Bounds().Dy(); y++ {
			yPost := y/16
			xPost := x/16
			postIndex := yPost * 16 + xPost;
			histogram[postIndex] += int(grayImage.GrayAt(x, y).Y)
		}
	}
	maxValue := 0
	for _, histValue := range histogram {
		maxValue = int(math.Max(float64(maxValue), float64(histValue)));
	}
	for index, hv := range histogram {
		histogram[index] = hv * 1000 / maxValue
	}
	return histogram
}

func (depot DepotImage) Save(path string) {
	io, err := os.Create(path)
	if (err != nil) {
		panic("cannot save")
	}
	png.Encode(io, depot.Image)
}

func (depot DepotImage) ItemDetectedAt(x int, y int, size int) DepotImage {
	var cropped DepotImage
	if (depot.Progress == EdgeDetection) {
		cropped = depot.Crop(x, y, size)
	} else {
		cropped = depot.GetEdge().Crop(x, y, size)
	}
	hist := cropped.getHistogram()
	confidence := getConfidence(GetBounds(), hist)
	fmt.Printf("%2f", confidence)
	// return true;
	return cropped;
}