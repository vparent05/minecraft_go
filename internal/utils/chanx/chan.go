package chanx

func TrySend[T any](ch chan<- T, t T) bool {
	select {
	case ch <- t:
		return true
	default:
		return false
	}
}

func TryReceive[T any](ch <-chan T) (T, bool) {
	select {
	case t := <-ch:
		return t, true
	default:
		return *new(T), false
	}
}
