// Package giotoast provides in-app toast notification widgets for Gio.
//
// A UI toast is a temporary, auto-dismissing notification that briefly appears on the screen to provide some type of feedback.
// It is called "toast" because typically it *pops-up* into view for a short while and then slides away once a timer ends — much like a piece of toast popping out of a toaster.
//
// Said another way, a toast is a brief message that appears at the bottom of the screen and auto-dismisses after a duration.
// It is similar to a Material Design "Snackbar".
//
// # Integration
//
// For most applications, use [Queue] — it handles multiple toasts gracefully,
// showing them one at a time:
//
//	type App struct {
//		toasts    giotoast.Queue
//		saveBtn   widget.Clickable
//		deleteBtn widget.Clickable
//	}
//
//	func (app *App) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
//		if app.saveBtn.Clicked(gtx) {
//			app.toasts.Enqueue("Profile saved", 3*time.Second, gtx.Now)
//		}
//
//		if app.deleteBtn.Clicked(gtx) {
//			app.toasts.EnqueueAction(giotoast.TypeError, "Message deleted", "UNDO", 5*time.Second, gtx.Now)
//		}
//
//		if app.toasts.ActionClicked(gtx) {
//			// undo the deletion
//		}
//
//		// overlay the toasts on top of your content
//		return layout.Stack{}.Layout(gtx,
//			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
//				return yourContent(gtx, th)
//			}),
//			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
//				return app.toasts.Layout(gtx, th)
//			}),
//		)
//	}
//
// The toast will fade in, remain visible for the given duration, and then fade out and disappear.
//
// # Toast Types
//
// Toasts can be categorized by type, each with a distinct background color and leading icon:
//
//   - [TypeNeutral] — dark background, no icon (the default)
//   - [TypeSuccess] — green background, checkmark icon
//   - [TypeError]   — red background, error icon
//   - [TypeWarning] — orange background, warning icon
//   - [TypeInfo]    — blue background, info icon
//
// Use [Queue.ShowType] to show a typed toast:
//
//	q.ShowType(giotoast.TypeSuccess, "Profile saved", 3*time.Second, gtx.Now)
//
// # Action Button
//
// Toasts can include an action button (for example, "UNDO"):
//
//	q.ShowAction(giotoast.TypeNeutral, "Message deleted", "UNDO", 5*time.Second, gtx.Now)
//
//	// Check if the action was clicked (before calling Layout):
//	if q.ActionClicked(gtx) {
//		// undo the deletion
//	}
//
// When the action button is clicked, the toast is automatically dismissed.
//
// # Close Button
//
// Every toast includes a close button (X icon) on the far right.
// Clicking it dismisses the toast with the same fade-out animation as auto-dismiss.
// No additional code is needed — the close button works automatically.
//
// You can also dismiss a toast programmatically:
//
//	q.Dismiss(gtx.Now)
//
// # Queue
//
// Use [Queue] to manage multiple toast messages, showing one at a time:
//
//	var q giotoast.Queue
//
//	// Enqueue multiple toasts:
//	q.Enqueue("First message", 3*time.Second, gtx.Now)
//	q.EnqueueType(giotoast.TypeSuccess, "Saved!", 3*time.Second, gtx.Now)
//	q.EnqueueAction(giotoast.TypeError, "Deleted", "UNDO", 5*time.Second, gtx.Now)
//
//	// Check for action clicks:
//	if q.ActionClicked(gtx) {
//		// handle the action
//	}
//
//	// In your layout:
//	layout.Stack{}.Layout(gtx,
//		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
//			return yourContent(gtx, th)
//		}),
//		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
//			return q.Layout(gtx, th)
//		}),
//	)
//
// When the first toast finishes (auto-dismiss or manual dismiss),
// the next pending toast is shown automatically.
//
// # Last-Write-Wins
//
// [Queue] also supports a "last-write-wins" pattern via [Queue.Show],
// [Queue.ShowType], and [Queue.ShowAction]. These methods replace the current
// toast immediately and clear any pending toasts:
//
//	// Replace whatever is showing (and discard anything queued):
//	q.Show("Connection lost", 5*time.Second, gtx.Now)
//
// Both patterns (FIFO and last-write-wins) can be mixed freely on the same [Queue].
//
// # Customization
//
// The default colors and durations can be changed by modifying the package-level variables:
//
//	giotoast.ColorBackground = color.NRGBA{R: 0x00, G: 0x00, B: 0x80, A: 0xFF}
//	giotoast.ColorText       = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
//	giotoast.ColorSuccess    = color.NRGBA{R: 0x00, G: 0x80, B: 0x00, A: 0xFF}
//	giotoast.ColorActionText = color.NRGBA{R: 0xFF, G: 0xD7, B: 0x00, A: 0xFF}
package giotoast
