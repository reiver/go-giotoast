package giotoast

import (
	"fmt"
	"image/color"

	"gioui.org/widget"

	"golang.org/x/exp/shiny/materialdesign/icons"
)

// Type represents the category of a toast notification.
//
// Each type has an associated background color and leading icon.
//
// See also:
//
//	• [TypeError]
//	• [TypeInfo]
//	• [TypeNeutral]
//	• [TypeSuccess]
//	• [TypeWarning]
type Type int

const (
	// TypeNeutral is the default toast type with no icon and the standard dark background.
	TypeNeutral Type = iota

	// TypeSuccess indicates a successful operation. Green background, checkmark icon.
	TypeSuccess

	// TypeError indicates a failure or error. Red background, error icon.
	TypeError

	// TypeWarning indicates a cautionary message. Orange background, warning icon.
	TypeWarning

	// TypeInfo indicates an informational message. Blue background, info icon.
	TypeInfo
)

// String returns the name of the toast type.
func (receiver Type) String() string {
	switch receiver {
	case TypeSuccess:
		return "Success"
	case TypeError:
		return "Error"
	case TypeWarning:
		return "Warning"
	case TypeInfo:
		return "Info"
	case TypeNeutral:
		return "Neutral"
	default:
		return fmt.Sprintf("Type(%d)", receiver)
	}
}

// Color returns the background color associated with the toast type.
func (receiver Type) Color() color.NRGBA {
	switch receiver {
	case TypeSuccess:
		return ColorSuccess
	case TypeError:
		return ColorError
	case TypeWarning:
		return ColorWarning
	case TypeInfo:
		return ColorInfo
	default:
		return ColorBackground
	}
}

// Icon returns the icon associated with the toast type, or nil for [TypeNeutral].
func (receiver Type) Icon() *widget.Icon {
	switch receiver {
	case TypeSuccess:
		return iconSuccess
	case TypeError:
		return iconError
	case TypeWarning:
		return iconWarning
	case TypeInfo:
		return iconInfo
	default:
		return nil
	}
}

var iconSuccess *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionCheckCircle)
	return icon
}()

var iconError *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.AlertError)
	return icon
}()

var iconWarning *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.AlertWarning)
	return icon
}()

var iconInfo *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionInfo)
	return icon
}()

var iconClose *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.NavigationClose)
	return icon
}()
