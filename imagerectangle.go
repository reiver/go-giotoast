package giotoast

import (
	"image"
)

func imageRectangle(width int, height int) image.Rectangle {
	return image.Rectangle{
		Max: image.Point{
			X: width,
			Y: height,
		},
	}
}
