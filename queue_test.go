package giotoast

import (
	"testing"

	"time"
)

func TestQueue_NilReceiver_Enqueue(t *testing.T) {
	var queue *Queue

	// should not panic, should return ErrReceiverNil
	var err error = queue.Enqueue("hello", 3*time.Second, time.Now())
	if ErrReceiverNil != err {
		t.Errorf("expected ErrReceiverNil on nil receiver, got %v", err)
	}
}

func TestQueue_NilReceiver_EnqueueType(t *testing.T) {
	var queue *Queue

	// should not panic, should return ErrReceiverNil
	var err error = queue.EnqueueType(TypeSuccess, "hello", 3*time.Second, time.Now())
	if ErrReceiverNil != err {
		t.Errorf("expected ErrReceiverNil on nil receiver, got %v", err)
	}
}

func TestQueue_NilReceiver_EnqueueAction(t *testing.T) {
	var queue *Queue

	// should not panic, should return ErrReceiverNil
	var err error = queue.EnqueueAction(TypeError, "deleted", "UNDO", 5*time.Second, time.Now())
	if ErrReceiverNil != err {
		t.Errorf("expected ErrReceiverNil on nil receiver, got %v", err)
	}
}

func TestQueue_NilReceiver_ActionClicked(t *testing.T) {
	var queue *Queue

	actual := queue.ActionClicked(layoutContext())
	if false != actual {
		t.Errorf("expected false on nil receiver, got %t", actual)
	}
}

func TestQueue_NilReceiver_Dismiss(t *testing.T) {
	var queue *Queue

	// should not panic
	queue.Dismiss(time.Now())
}

func TestQueue_NilReceiver_Layout(t *testing.T) {
	var queue *Queue

	// should not panic, should return zero dimensions
	var dims = queue.Layout(layoutContext(), nil)
	if 0 != dims.Size.X || 0 != dims.Size.Y {
		t.Errorf("expected zero dimensions on nil receiver, got %v", dims.Size)
	}
}

func TestQueue_Enqueue_ShowsImmediately(t *testing.T) {
	var queue Queue

	var now time.Time = time.Now()

	queue.Enqueue("first", 3*time.Second, now)

	if !queue.current.Visible() {
		t.Error("expected first toast to be visible immediately")
	}

	if "first" != queue.current.message {
		t.Errorf("expected message %q, got %q", "first", queue.current.message)
	}
}

func TestQueue_EnqueueType_ShowsImmediately(t *testing.T) {
	var queue Queue

	var now time.Time = time.Now()

	queue.EnqueueType(TypeSuccess, "saved", 3*time.Second, now)

	if !queue.current.Visible() {
		t.Error("expected first toast to be visible immediately")
	}

	if TypeSuccess != queue.current.toastType {
		t.Errorf("expected TypeSuccess, got %d", queue.current.toastType)
	}
}

func TestQueue_EnqueueAction_ShowsImmediately(t *testing.T) {
	var queue Queue

	var now time.Time = time.Now()

	queue.EnqueueAction(TypeError, "deleted", "UNDO", 5*time.Second, now)

	if !queue.current.Visible() {
		t.Error("expected first toast to be visible immediately")
	}

	if "UNDO" != queue.current.action {
		t.Errorf("expected action %q, got %q", "UNDO", queue.current.action)
	}
}

func TestQueue_Enqueue_QueuesSecond(t *testing.T) {
	var queue Queue

	var now time.Time = time.Now()

	queue.Enqueue("first", 3*time.Second, now)
	queue.Enqueue("second", 3*time.Second, now)

	if "first" != queue.current.message {
		t.Errorf("expected current message %q, got %q", "first", queue.current.message)
	}

	if 1 != len(queue.pending) {
		t.Errorf("expected 1 pending toast, got %d", len(queue.pending))
	}

	if "second" != queue.pending[0].message {
		t.Errorf("expected pending message %q, got %q", "second", queue.pending[0].message)
	}
}

func TestQueue_Enqueue_DefaultDuration(t *testing.T) {
	var queue Queue

	var now time.Time = time.Now()

	queue.Enqueue("hello", 0, now)

	if DefaultDuration != queue.current.duration {
		t.Errorf("expected default duration %v, got %v", DefaultDuration, queue.current.duration)
	}
}

func TestQueue_EnqueueAction_DefaultDuration(t *testing.T) {
	var queue Queue

	var now time.Time = time.Now()

	queue.EnqueueAction(TypeError, "deleted", "UNDO", 0, now)

	if DefaultActionDuration != queue.current.duration {
		t.Errorf("expected duration %v, got %v", DefaultActionDuration, queue.current.duration)
	}
}

func TestQueue_Show(t *testing.T) {
	var queue Queue

	var now time.Time = time.Now()

	queue.Show("hello", 3*time.Second, now)

	if !queue.current.Visible() {
		t.Error("expected toast to be visible")
	}

	if "hello" != queue.current.message {
		t.Errorf("expected message %q, got %q", "hello", queue.current.message)
	}

	if TypeNeutral != queue.current.toastType {
		t.Errorf("expected TypeNeutral, got %v", queue.current.toastType)
	}
}

func TestQueue_ShowType(t *testing.T) {
	var queue Queue

	var now time.Time = time.Now()

	queue.ShowType(TypeSuccess, "saved", 3*time.Second, now)

	if !queue.current.Visible() {
		t.Error("expected toast to be visible")
	}

	if TypeSuccess != queue.current.toastType {
		t.Errorf("expected TypeSuccess, got %v", queue.current.toastType)
	}
}

func TestQueue_ShowAction(t *testing.T) {
	var queue Queue

	var now time.Time = time.Now()

	queue.ShowAction(TypeError, "deleted", "UNDO", 5*time.Second, now)

	if !queue.current.Visible() {
		t.Error("expected toast to be visible")
	}

	if "UNDO" != queue.current.action {
		t.Errorf("expected action %q, got %q", "UNDO", queue.current.action)
	}
}

func TestQueue_Show_ClearsPending(t *testing.T) {
	var queue Queue

	var now time.Time = time.Now()

	// enqueue several toasts
	queue.Enqueue("first", 3*time.Second, now)
	queue.Enqueue("second", 3*time.Second, now)
	queue.Enqueue("third", 3*time.Second, now)

	if 2 != len(queue.pending) {
		t.Fatalf("expected 2 pending toasts, got %d", len(queue.pending))
	}

	// Show should clear pending and replace current
	queue.Show("urgent", 3*time.Second, now)

	if 0 != len(queue.pending) {
		t.Errorf("expected 0 pending toasts after Show, got %d", len(queue.pending))
	}

	if "urgent" != queue.current.message {
		t.Errorf("expected message %q, got %q", "urgent", queue.current.message)
	}
}

func TestQueue_Show_DefaultDuration(t *testing.T) {
	var queue Queue

	var now time.Time = time.Now()

	queue.Show("hello", 0, now)

	if DefaultDuration != queue.current.duration {
		t.Errorf("expected default duration %v, got %v", DefaultDuration, queue.current.duration)
	}
}

func TestQueue_ShowAction_DefaultDuration(t *testing.T) {
	var queue Queue

	var now time.Time = time.Now()

	queue.ShowAction(TypeError, "deleted", "UNDO", 0, now)

	if DefaultActionDuration != queue.current.duration {
		t.Errorf("expected duration %v, got %v", DefaultActionDuration, queue.current.duration)
	}
}

func TestQueue_NilReceiver_Show(t *testing.T) {
	var queue *Queue

	// should not panic
	queue.Show("hello", 3*time.Second, time.Now())
}

func TestQueue_NilReceiver_ShowType(t *testing.T) {
	var queue *Queue

	// should not panic
	queue.ShowType(TypeSuccess, "hello", 3*time.Second, time.Now())
}

func TestQueue_NilReceiver_ShowAction(t *testing.T) {
	var queue *Queue

	// should not panic
	queue.ShowAction(TypeError, "deleted", "UNDO", 5*time.Second, time.Now())
}

func TestQueue_Enqueue_QueueFull(t *testing.T) {
	var queue Queue

	var now time.Time = time.Now()

	// show the first toast so subsequent enqueues go to pending
	queue.Enqueue("first", 3*time.Second, now)

	// fill up the pending queue
	for i := 0; i < MaxQueueSize; i++ {
		var err error = queue.Enqueue("pending", 3*time.Second, now)
		if nil != err {
			t.Fatalf("expected nil error on enqueue %d, got %v", i, err)
		}
	}

	if MaxQueueSize != len(queue.pending) {
		t.Fatalf("expected %d pending toasts, got %d", MaxQueueSize, len(queue.pending))
	}

	// next enqueue should fail
	var err error = queue.Enqueue("overflow", 3*time.Second, now)
	if ErrQueueFull != err {
		t.Errorf("expected ErrQueueFull, got %v", err)
	}
}
