package giotoast

import (
	"image"
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"gioui.org/x/component"
)

// Toast represents an in-app toast notification that appears at the bottom
// of the screen and auto-dismisses after a duration.
//
// A Toast is similar to a Material Design "Snackbar".
//
// Example usage:
//
//	var t giotoast.Toast
//
//	// Show a plain toast:
//	t.Show("Profile saved", 3*time.Second, gtx.Now)
//
//	// Show a typed toast with an icon:
//	t.ShowType(giotoast.TypeSuccess, "Profile saved", 3*time.Second, gtx.Now)
//
//	// Show a toast with an action button:
//	t.ShowAction(giotoast.TypeNeutral, "Message deleted", "UNDO", 5*time.Second, gtx.Now)
//
//	// Check if the action button was clicked:
//	if t.ActionClicked(gtx) {
//		// handle undo
//	}
//
//	// In your layout, overlay the toast on top of your content:
//	layout.Stack{}.Layout(gtx,
//		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
//			return yourContent(gtx, th)
//		}),
//		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
//			return t.Layout(gtx, th)
//		}),
//	)
type Toast struct {
	// shown guards against VisibilityAnimation's zero value being Visible (iota=0).
	// Without this, a zero-value Toast would report as visible.
	shown           bool
	message         string
	toastType       Type
	action          string
	actionBtn       widget.Clickable
	closeBtn        widget.Clickable
	showTime        time.Time
	duration        time.Duration
	anim            component.VisibilityAnimation
	bgColor         color.NRGBA
	textColor       color.NRGBA
	iconTint        color.NRGBA
	actionTextColor color.NRGBA
}

// Show shows a [TypeNeutral] toast with the given message that will auto-dismiss after the given duration.
//
// If duration is less-than or equal-to zero, a default duration of 3 seconds is used.
func (receiver *Toast) Show(message string, duration time.Duration, now time.Time) {
	if nil == receiver {
		return
	}

	receiver.show(TypeNeutral, message, "", duration, now)
}

// ShowType shows a typed toast with the given message that will auto-dismiss after the given duration.
//
// The toast type determines the background color and leading icon.
//
// If duration is less-than or equal-to zero, a default duration of 3 seconds is used.
func (receiver *Toast) ShowType(toastType Type, message string, duration time.Duration, now time.Time) {
	if nil == receiver {
		return
	}

	receiver.show(toastType, message, "", duration, now)
}

// ShowAction shows a toast with a type, message, and an action button.
//
// Use [Toast.ActionClicked] to check if the action button was clicked.
//
// If duration is less-than or equal-to zero, a default duration of 5 seconds is used
// (longer than the default for non-action toasts, to give the user time to read and act).
func (receiver *Toast) ShowAction(toastType Type, message string, action string, duration time.Duration, now time.Time) {
	if nil == receiver {
		return
	}

	if duration <= 0 {
		duration = DefaultActionDuration
	}

	receiver.show(toastType, message, action, duration, now)
}

// ActionClicked reports whether the action button was clicked since the last call to ActionClicked.
//
// When the action is clicked, the toast is also dismissed.
func (receiver *Toast) ActionClicked(gtx layout.Context) bool {
	if nil == receiver {
		return false
	}

	if !receiver.shown {
		return false
	}

	if !receiver.anim.Visible() {
		return false
	}

	if receiver.actionBtn.Clicked(gtx) {
		receiver.anim.Disappear(gtx.Now)
		return true
	}

	return false
}

// Dismiss dismisses the toast immediately (with a fade-out animation).
func (receiver *Toast) Dismiss(now time.Time) {
	if nil == receiver {
		return
	}

	receiver.anim.Disappear(now)
}

// Visible reports whether the toast is currently visible (including while animating).
func (receiver *Toast) Visible() bool {
	if nil == receiver {
		return false
	}

	if !receiver.shown {
		return false
	}

	return receiver.anim.Visible()
}

// Layout draws the toast. It should be called as an overlay on top of your main content.
//
// If the toast is not visible, Layout is a no-op and returns zero dimensions.
func (receiver *Toast) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if nil == receiver {
		return layout.Dimensions{}
	}

	if !receiver.shown || !receiver.anim.Visible() {
		return layout.Dimensions{}
	}

	// close button dismiss
	if receiver.closeBtn.Clicked(gtx) {
		receiver.anim.Disappear(gtx.Now)
	}

	// auto-dismiss
	if 0 < receiver.duration && gtx.Now.After(receiver.showTime.Add(receiver.duration)) {
		receiver.anim.Disappear(gtx.Now)
	}

	// schedule a redraw at dismissal time so the fade-out triggers even if nothing else causes a frame
	if receiver.anim.Visible() && !receiver.anim.Animating() {
		var dismissAt time.Time = receiver.showTime.Add(receiver.duration)
		if gtx.Now.Before(dismissAt) {
			gtx.Execute(op.InvalidateCmd{At: dismissAt})
		}
	}

	var revealed float32 = receiver.anim.Revealed(gtx)

	return layout.S.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Bottom: unit.Dp(24),
			Left:   unit.Dp(16),
			Right:  unit.Dp(16),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return receiver.layoutCard(gtx, th, revealed)
		})
	})
}

func (receiver *Toast) show(toastType Type, message string, action string, duration time.Duration, now time.Time) {
	if duration <= 0 {
		duration = DefaultDuration
	}

	receiver.shown = true
	receiver.toastType = toastType
	receiver.message = message
	receiver.action = action
	receiver.duration = duration
	receiver.showTime = now
	receiver.bgColor = toastType.Color()
	receiver.textColor = ColorText
	receiver.iconTint = ColorIconTint
	receiver.actionTextColor = ColorActionText
	receiver.anim.Duration = DefaultAnimationDuration
	receiver.anim.State = component.Invisible
	receiver.anim.Appear(now)
}

func (receiver *Toast) layoutCard(gtx layout.Context, th *material.Theme, revealed float32) layout.Dimensions {

	var maxWidth int = gtx.Dp(unit.Dp(568))
	if gtx.Constraints.Max.X > maxWidth {
		gtx.Constraints.Max.X = maxWidth
	}

	var macro op.MacroOp = op.Record(gtx.Ops)

	dims := layout.Inset{
		Top:    unit.Dp(12),
		Bottom: unit.Dp(12),
		Left:   unit.Dp(16),
		Right:  unit.Dp(16),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Alignment: layout.Middle,
		}.Layout(gtx,
			receiver.layoutIcon(gtx),
			receiver.layoutMessage(gtx, th),
			receiver.layoutAction(gtx, th),
			receiver.layoutClose(gtx, th),
		)
	})

	var call op.CallOp = macro.Stop()

	// background
	var rr int = gtx.Dp(unit.Dp(4))
	var rrect clip.RRect = clip.RRect{
		Rect: imageRectangle(dims.Size.X, dims.Size.Y),
		SE:   rr,
		SW:   rr,
		NE:   rr,
		NW:   rr,
	}

	// fade the entire card (background + content) together
	var opacityStack paint.OpacityStack = paint.PushOpacity(gtx.Ops, revealed)

	var bgColor color.NRGBA = receiver.bgColor

	paintStack := rrect.Push(gtx.Ops)
	paint.ColorOp{Color: bgColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	paintStack.Pop()

	// content
	contentStack := rrect.Push(gtx.Ops)
	call.Add(gtx.Ops)
	contentStack.Pop()

	opacityStack.Pop()

	return dims
}

func (receiver *Toast) layoutIcon(gtx layout.Context) layout.FlexChild {
	var icon *widget.Icon = receiver.toastType.Icon()
	if nil == icon {
		return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{}
		})
	}

	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Right: unit.Dp(12),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			var sizePx int = gtx.Dp(unit.Dp(20))
			gtx.Constraints.Min = image.Point{X: sizePx, Y: sizePx}
			gtx.Constraints.Max = image.Point{X: sizePx, Y: sizePx}
			return icon.Layout(gtx, receiver.iconTint)
		})
	})
}

func (receiver *Toast) layoutMessage(gtx layout.Context, th *material.Theme) layout.FlexChild {
	return layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		var lbl material.LabelStyle = material.Body2(th, receiver.message)
		lbl.Color = receiver.textColor
		return lbl.Layout(gtx)
	})
}

func (receiver *Toast) layoutAction(gtx layout.Context, th *material.Theme) layout.FlexChild {
	if "" == receiver.action {
		return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{}
		})
	}

	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Left: unit.Dp(8),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			var btn material.ButtonStyle = material.Button(th, &receiver.actionBtn, receiver.action)
			btn.Color = receiver.actionTextColor
			btn.Background = color.NRGBA{}
			btn.Inset = layout.Inset{
				Top:    unit.Dp(8),
				Bottom: unit.Dp(8),
				Left:   unit.Dp(8),
				Right:  unit.Dp(8),
			}
			return btn.Layout(gtx)
		})
	})
}

func (receiver *Toast) layoutClose(gtx layout.Context, th *material.Theme) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Left: unit.Dp(4),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			var btn material.IconButtonStyle = material.IconButton(th, &receiver.closeBtn, iconClose, "Dismiss")
			btn.Color = receiver.textColor
			btn.Background = color.NRGBA{}
			btn.Size = unit.Dp(18)
			btn.Inset = layout.UniformInset(unit.Dp(4))
			return btn.Layout(gtx)
		})
	})
}
