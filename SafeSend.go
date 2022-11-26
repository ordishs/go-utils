package utils

// send returns true if it was able to send t on channel c.
// It returns false if c is closed.
func SafeSend[T any](ch chan T, t T) bool {
	defer func() {
		_ = recover()
	}()

	select {
	case <-ch:
		return false // Channel is closed
	default:
		ch <- t
	}

	return true
}
