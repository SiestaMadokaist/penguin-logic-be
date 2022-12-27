package DepotRecognition

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

type Cropable interface {
	SubImage(r image.Rectangle) image.Image
}

func GetBounds() [][2]int {
	detector, err := os.Open("./features/detector.txt")
	if (err != nil) {
		panic(errors.New("cannot open fetaures/detector"))
	}
	b, err := io.ReadAll(detector)
	if (err != nil) {
		panic(errors.New("cannot open features/detector"))
	}
	boundaries := make([][2]int, 256)
	content := strings.Split(string(b), "\n");
	for i, line := range content {
		v := strings.Split(line, " ")
		low, _ := strconv.ParseInt(v[0], 10, 16)
		high, _ := strconv.ParseInt(v[1], 10, 16)
		boundaries[i] = [2]int{ int(low), int(high) }
	}
	return boundaries;
}

type Rejection struct {
	X int;
	Y int;
	Bound [2]int;
	V int;
}
func getRejectedLocation(bounds [][2]int, values []int) []Rejection {
	inRange := 0
	var rejected []Rejection;
	for i, bound := range bounds {
		lowerbound := bound[0]
		upperbound := bound[1]
		value := values[i];
		X := i / 16
		Y := i % 16
		v := Rejection{ X: X, Y: Y, Bound: bound, V: value }
		if lowerbound <= value && value <= upperbound {
			inRange++
		} else {
			rejected = append(rejected, v)
		}
	}
	return rejected;
}

func (depot DepotImage) Crop(x int, y int, size int) DepotImage {
	original, ok := depot.Image.(Cropable)
	if (!ok) { panic("cannot convert image to cropable"); }
	cropped := original.SubImage(image.Rect(x, y, x + size, y + size));
	return DepotImage{Image: cropped, Progress: depot.Progress }
}

func (depot DepotImage) getHistogram() []int {
	histogram := make([]int, 256);
	grayImage := depot.Image.(*image.Gray);
	min := grayImage.Bounds().Min
	for x := 0; x < grayImage.Bounds().Dx(); x++ {
		for y := 0; y < grayImage.Bounds().Dy(); y++ {
			yPost := y/16
			xPost := x/16
			postIndex := yPost * 16 + xPost;
			xLoc := x + min.X
			yLoc := y + min.Y
			grayValue := grayImage.GrayAt(xLoc, yLoc)
			histogram[postIndex] += int(grayValue.Y)
		}
	}
	maxValue := 1
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

func (cropped DepotImage) Rejections() (rejection []Rejection) {
	hist := cropped.getHistogram()
	rejection = getRejectedLocation(GetBounds(), hist)
	return rejection;
}

func (depot DepotImage) ItemDetectedAt(x int, y int, size int) (cropped DepotImage, confidence float32) {
	if (depot.Progress == EdgeDetection) {
		cropped = depot.Crop(x, y, size)
	} else {
		cropped = depot.GetEdge().Crop(x, y, size)
	}
	rejection := cropped.Rejections()
	score := float32(256 - len(rejection))
	confidence = score / 256;
	fmt.Printf("%2f %d %2f \n", score, len(rejection), confidence);
	// confidence = float32(256 - len(rejection)) / 256
	return depot.Crop(x, y, size), confidence;
}