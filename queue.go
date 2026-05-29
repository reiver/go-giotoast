package giotoast

import (
	"time"

	"gioui.org/layout"
	"gioui.org/widget/material"
)

// Queue manages a queue of toast messages, showing one at a time.
//
// When a toast is dismissed (either by auto-dismiss or manual dismiss),
// the next toast in the queue is shown.
//
// Example usage:
//
//	var q giotoast.Queue
//
//	// Enqueue toasts:
//	q.Enqueue("First message", 3*time.Second, gtx.Now)
//	q.EnqueueType(giotoast.TypeSuccess, "Saved!", 3*time.Second, gtx.Now)
//	q.EnqueueAction(giotoast.TypeError, "Deleted", "UNDO", 5*time.Second, gtx.Now)
//
//	// Check for action clicks:
//	if q.ActionClicked(gtx) {
//		// handle the action
//	}
//
//	// In your layout, overlay the queue on top of your content:
//	layout.Stack{}.Layout(gtx,
//		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
//			return yourContent(gtx, th)
//		}),
//		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
//			return q.Layout(gtx, th)
//		}),
//	)
type Queue struct {
	current Toast
	pending []pendingToast
}

type pendingToast struct {
	toastType Type
	message   string
	action    string
	duration  time.Duration
}

// Enqueue adds a [TypeNeutral] toast message to the queue.
//
// If no toast is currently visible, the message is shown immediately.
// Otherwise, it is queued and shown after the current toast dismisses.
//
// If duration is less-than or equal-to zero, a default duration of 3 seconds is used.
//
// Enqueue returns [ErrQueueFull] if the queue has reached [MaxQueueSize] pending toasts.
func (receiver *Queue) Enqueue(message string, duration time.Duration, now time.Time) error {
	if nil == receiver {
		return ErrReceiverNil
	}

	return receiver.enqueue(TypeNeutral, message, "", duration, now)
}

// EnqueueType adds a typed toast message to the queue.
//
// If no toast is currently visible, the message is shown immediately.
// Otherwise, it is queued and shown after the current toast dismisses.
//
// If duration is less-than or equal-to zero, a default duration of 3 seconds is used.
//
// EnqueueType returns [ErrQueueFull] if the queue has reached [MaxQueueSize] pending toasts.
func (receiver *Queue) EnqueueType(toastType Type, message string, duration time.Duration, now time.Time) error {
	if nil == receiver {
		return ErrReceiverNil
	}

	return receiver.enqueue(toastType, message, "", duration, now)
}

// EnqueueAction adds a typed toast with an action button to the queue.
//
// Use [Queue.ActionClicked] to check if the action button was clicked.
//
// If no toast is currently visible, the message is shown immediately.
// Otherwise, it is queued and shown after the current toast dismisses.
//
// If duration is less-than or equal-to zero, a default duration of 5 seconds is used.
//
// EnqueueAction returns [ErrQueueFull] if the queue has reached [MaxQueueSize] pending toasts.
func (receiver *Queue) EnqueueAction(toastType Type, message string, action string, duration time.Duration, now time.Time) error {
	if nil == receiver {
		return ErrReceiverNil
	}

	return receiver.enqueue(toastType, message, action, duration, now)
}

// ActionClicked reports whether the action button on the current toast was clicked.
//
// When the action is clicked, the current toast is also dismissed.
func (receiver *Queue) ActionClicked(gtx layout.Context) bool {
	if nil == receiver {
		return false
	}

	return receiver.current.ActionClicked(gtx)
}

// Dismiss dismisses the current toast immediately (with a fade-out animation).
func (receiver *Queue) Dismiss(now time.Time) {
	if nil == receiver {
		return
	}

	receiver.current.Dismiss(now)
}

// Layout draws the current toast. It should be called as an overlay on top of your main content.
//
// Layout also advances the queue — when the current toast finishes disappearing,
// the next pending toast is shown.
func (receiver *Queue) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if nil == receiver {
		return layout.Dimensions{}
	}

	// advance the queue if the current toast has finished
	if !receiver.current.Visible() && 0 < len(receiver.pending) {
		var next pendingToast = receiver.pending[0]
		// Zero the consumed slot so its strings can be GC'd while the backing array is still alive.
		receiver.pending[0] = pendingToast{}
		receiver.pending = receiver.pending[1:]

		if "" != next.action {
			receiver.current.ShowAction(next.toastType, next.message, next.action, next.duration, gtx.Now)
		} else {
			receiver.current.ShowType(next.toastType, next.message, next.duration, gtx.Now)
		}
	}

	return receiver.current.Layout(gtx, th)
}

func (receiver *Queue) enqueue(toastType Type, message string, action string, duration time.Duration, now time.Time) error {
	if duration <= 0 {
		if "" != action {
			duration = DefaultActionDuration
		} else {
			duration = DefaultDuration
		}
	}

	if !receiver.current.Visible() {
		if "" != action {
			receiver.current.ShowAction(toastType, message, action, duration, now)
		} else {
			receiver.current.ShowType(toastType, message, duration, now)
		}
		return nil
	}

	if len(receiver.pending) >= MaxQueueSize {
		return ErrQueueFull
	}

	receiver.pending = append(receiver.pending, pendingToast{
		toastType: toastType,
		message:   message,
		action:    action,
		duration:  duration,
	})

	return nil
}
