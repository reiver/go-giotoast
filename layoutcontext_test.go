package giotoast

import (
	"image"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
)

func layoutContext() layout.Context {
	var ops op.Ops
	return layout.Context{
		Ops: &ops,
		Now: time.Now(),
		Constraints: layout.Constraints{
			Max: image.Point{X: 1080, Y: 1920},
		},
	}
}
