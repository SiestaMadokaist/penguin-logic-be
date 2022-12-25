package DepotRecognition

import "image"

type Progress int64

const (
	Original Progress = iota
	EdgeDetected
)

type DepotImage struct {
	Image image.Image
	Progress Progress
}

func NewDepot(image image.Image) DepotImage {
	return DepotImage{Image: image, Progress: Original}
}