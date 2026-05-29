package giotoast

import (
	"testing"

	"time"

	"gioui.org/widget/material"
)

func TestToast_NilReceiver_Show(t *testing.T) {
	var toast *Toast

	// should not panic
	toast.Show("hello", 3*time.Second, time.Now())
}

func TestToast_NilReceiver_ShowType(t *testing.T) {
	var toast *Toast

	// should not panic
	toast.ShowType(TypeSuccess, "hello", 3*time.Second, time.Now())
}

func TestToast_NilReceiver_ShowAction(t *testing.T) {
	var toast *Toast

	// should not panic
	toast.ShowAction(TypeError, "deleted", "UNDO", 5*time.Second, time.Now())
}

func TestToast_NilReceiver_ActionClicked(t *testing.T) {
	var toast *Toast

	actual := toast.ActionClicked(layoutContext())
	if false != actual {
		t.Errorf("expected false on nil receiver, got %t", actual)
	}
}

func TestToast_NilReceiver_Dismiss(t *testing.T) {
	var toast *Toast

	// should not panic
	toast.Dismiss(time.Now())
}

func TestToast_NilReceiver_Visible(t *testing.T) {
	var toast *Toast

	actual := toast.Visible()
	if false != actual {
		t.Errorf("expected false on nil receiver, got %t", actual)
	}
}

func TestToast_NilReceiver_Layout(t *testing.T) {
	var toast *Toast

	// should not panic, should return zero dimensions
	var dims = toast.Layout(layoutContext(), nil)
	if 0 != dims.Size.X || 0 != dims.Size.Y {
		t.Errorf("expected zero dimensions on nil receiver, got %v", dims.Size)
	}
}

func TestToast_NotVisible_Initially(t *testing.T) {
	var toast Toast

	if toast.Visible() {
		t.Error("expected toast to not be visible initially")
	}
}

func TestToast_Show(t *testing.T) {
	var toast Toast

	var now time.Time = time.Now()

	toast.Show("hello", 3*time.Second, now)

	if !toast.Visible() {
		t.Error("expected toast to be visible after Show()")
	}

	if "hello" != toast.message {
		t.Errorf("expected message %q, got %q", "hello", toast.message)
	}

	if TypeNeutral != toast.toastType {
		t.Errorf("expected TypeNeutral, got %d", toast.toastType)
	}
}

func TestToast_Show_DefaultDuration(t *testing.T) {
	var toast Toast

	var now time.Time = time.Now()

	toast.Show("hello", 0, now)

	if DefaultDuration != toast.duration {
		t.Errorf("expected default duration %v, got %v", DefaultDuration, toast.duration)
	}
}

func TestToast_ShowType(t *testing.T) {
	var toast Toast

	var now time.Time = time.Now()

	toast.ShowType(TypeSuccess, "saved", 3*time.Second, now)

	if !toast.Visible() {
		t.Error("expected toast to be visible after ShowType()")
	}

	if TypeSuccess != toast.toastType {
		t.Errorf("expected TypeSuccess, got %d", toast.toastType)
	}

	if "saved" != toast.message {
		t.Errorf("expected message %q, got %q", "saved", toast.message)
	}
}

func TestToast_ShowAction(t *testing.T) {
	var toast Toast

	var now time.Time = time.Now()

	toast.ShowAction(TypeError, "deleted", "UNDO", 5*time.Second, now)

	if !toast.Visible() {
		t.Error("expected toast to be visible after ShowAction()")
	}

	if TypeError != toast.toastType {
		t.Errorf("expected TypeError, got %d", toast.toastType)
	}

	if "deleted" != toast.message {
		t.Errorf("expected message %q, got %q", "deleted", toast.message)
	}

	if "UNDO" != toast.action {
		t.Errorf("expected action %q, got %q", "UNDO", toast.action)
	}
}

func TestToast_ShowAction_DefaultDuration(t *testing.T) {
	var toast Toast

	var now time.Time = time.Now()

	toast.ShowAction(TypeError, "deleted", "UNDO", 0, now)

	if DefaultActionDuration != toast.duration {
		t.Errorf("expected duration %v, got %v", DefaultActionDuration, toast.duration)
	}
}

func TestToast_Dismiss(t *testing.T) {
	var toast Toast

	var now time.Time = time.Now()

	toast.Show("hello", 3*time.Second, now)
	toast.Dismiss(now.Add(1 * time.Second))

	// after Dismiss, the toast may still be visible while animating,
	// but it should no longer be in the "appeared" state
	// (it will be animating towards invisible)
}

func TestToast_Layout_Visible(t *testing.T) {
	var toast Toast
	var th *material.Theme = material.NewTheme()
	var gtx = layoutContext()

	toast.Show("hello", 3*time.Second, gtx.Now)

	var dims = toast.Layout(gtx, th)
	if 0 == dims.Size.X || 0 == dims.Size.Y {
		t.Errorf("expected non-zero dimensions for visible toast, got %v", dims.Size)
	}
}

func TestToast_Layout_VisibleWithAction(t *testing.T) {
	var toast Toast
	var th *material.Theme = material.NewTheme()
	var gtx = layoutContext()

	toast.ShowAction(TypeError, "deleted", "UNDO", 5*time.Second, gtx.Now)

	var dims = toast.Layout(gtx, th)
	if 0 == dims.Size.X || 0 == dims.Size.Y {
		t.Errorf("expected non-zero dimensions for visible action toast, got %v", dims.Size)
	}
}
