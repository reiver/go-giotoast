package giotoast

import (
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
)

// Queue manages a queue of toast messages, showing one at a time.
//
// Like most Gio widgets, a Queue is not safe for concurrent use.
// All methods must be called from the Gio main goroutine (the event loop).
//
// Queue supports two usage patterns:
//
//   - FIFO (first-in, first-out): Use [Queue.Enqueue], [Queue.EnqueueType], and
//     [Queue.EnqueueAction] to add toasts to the queue. Each toast is shown in order
//     after the previous one dismisses.
//
//   - Last-write-wins: Use [Queue.Show], [Queue.ShowType], and [Queue.ShowAction] to
//     replace the current toast immediately and clear any pending toasts.
//
// Both patterns can be mixed freely on the same Queue.
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
	current toast
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

// Show shows a [TypeNeutral] toast immediately, clearing any pending toasts.
// This is the "last-write-wins" counterpart to [Queue.Enqueue].
//
// If a toast is already visible, it is replaced immediately (without a fade-out).
//
// If duration is less-than or equal-to zero, a default duration of 3 seconds is used.
func (receiver *Queue) Show(message string, duration time.Duration, now time.Time) {
	if nil == receiver {
		return
	}

	receiver.pending = nil
	receiver.current.Show(message, duration, now)
}

// ShowType shows a typed toast immediately, clearing any pending toasts.
// This is the "last-write-wins" counterpart to [Queue.EnqueueType].
//
// If a toast is already visible, it is replaced immediately (without a fade-out).
//
// If duration is less-than or equal-to zero, a default duration of 3 seconds is used.
func (receiver *Queue) ShowType(toastType Type, message string, duration time.Duration, now time.Time) {
	if nil == receiver {
		return
	}

	receiver.pending = nil
	receiver.current.ShowType(toastType, message, duration, now)
}

// ShowAction shows a toast with an action button immediately, clearing any pending toasts.
// This is the "last-write-wins" counterpart to [Queue.EnqueueAction].
//
// Use [Queue.ActionClicked] to check if the action button was clicked.
//
// If a toast is already visible, it is replaced immediately (without a fade-out).
//
// If duration is less-than or equal-to zero, a default duration of 5 seconds is used.
func (receiver *Queue) ShowAction(toastType Type, message string, action string, duration time.Duration, now time.Time) {
	if nil == receiver {
		return
	}

	receiver.pending = nil
	receiver.current.ShowAction(toastType, message, action, duration, now)
}

// ActionClicked reports whether the action button on the current toast was clicked.
//
// When the action is clicked, the current toast is also dismissed.
//
// ActionClicked does not identify which toast's action was clicked — it applies
// to whichever toast is currently showing. If you need different handlers for
// different actions, use [Queue.ShowAction] (last-write-wins) so that only one
// action toast is active at a time.
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

// Visible reports whether a toast is currently visible (including while animating).
func (receiver *Queue) Visible() bool {
	if nil == receiver {
		return false
	}

	return receiver.current.Visible()
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

	// Ensure a frame fires after the current toast's fade-out completes so the queue can advance.
	if receiver.current.Visible() && 0 < len(receiver.pending) {
		gtx.Execute(op.InvalidateCmd{})
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
