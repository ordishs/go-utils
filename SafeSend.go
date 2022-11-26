package utils

// SafeSend is a helper function that sends a value on a channel and returns.
// It is safe to use with closed channels.
// It uses go generics to allow for any type of channel.
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
