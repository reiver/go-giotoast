package giotoast

import (
	"testing"
)

func TestQueue_NilReceiver_Visible(t *testing.T) {
	var queue *Queue

	actual := queue.Visible()
	if false != actual {
		t.Errorf("expected false on nil receiver, got %t", actual)
	}
}
