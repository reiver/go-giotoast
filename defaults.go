package giotoast

import (
	"image/color"
	"time"
)

// MaxQueueSize is the maximum number of pending toasts allowed in a [Queue].
// Attempts to enqueue beyond this limit will return [ErrQueueFull].
const MaxQueueSize int = 1024

const (
	// DefaultDuration is the default duration a toast is visible before auto-dismissing.
	DefaultDuration time.Duration = 3 * time.Second

	// DefaultActionDuration is the default duration a toast with an action button is visible before auto-dismissing.
	// This is longer than [DefaultDuration] to give the user time to read and act.
	DefaultActionDuration time.Duration = 5 * time.Second

	// DefaultAnimationDuration is the default duration of the fade-in and fade-out animation.
	DefaultAnimationDuration time.Duration = 250 * time.Millisecond
)

var (
	// ColorBackground is the default background color for [TypeNeutral] toasts.
	// This is the Material Design dark surface color.
	ColorBackground color.NRGBA = color.NRGBA{R: 0x32, G: 0x32, B: 0x32, A: 0xFF}

	// ColorText is the default text color of the toast.
	ColorText color.NRGBA = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

	// ColorIconTint is the color used to tint the leading icon.
	ColorIconTint color.NRGBA = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

	// ColorActionText is the color used for the action button text.
	ColorActionText color.NRGBA = color.NRGBA{R: 0xBB, G: 0x86, B: 0xFC, A: 0xFF}

	// ColorSuccess is the background color for [TypeSuccess] toasts.
	ColorSuccess color.NRGBA = color.NRGBA{R: 0x2E, G: 0x7D, B: 0x32, A: 0xFF}

	// ColorError is the background color for [TypeError] toasts.
	ColorError color.NRGBA = color.NRGBA{R: 0xC6, G: 0x28, B: 0x28, A: 0xFF}

	// ColorWarning is the background color for [TypeWarning] toasts.
	ColorWarning color.NRGBA = color.NRGBA{R: 0xE6, G: 0x5C, B: 0x00, A: 0xFF}

	// ColorInfo is the background color for [TypeInfo] toasts.
	ColorInfo color.NRGBA = color.NRGBA{R: 0x15, G: 0x65, B: 0xC0, A: 0xFF}
)
