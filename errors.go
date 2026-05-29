package giotoast

import (
	"codeberg.org/reiver/go-erorr"
)

const (
	// ErrQueueFull is returned when a [Queue] enqueue operation fails because the queue has reached [MaxQueueSize] pending toasts.
	ErrQueueFull = erorr.Error("giotoast: queue is full")

	// ErrReceiverNil is returned when a receiver is nil.
	ErrReceiverNil = erorr.Error("nil receiver")
)
