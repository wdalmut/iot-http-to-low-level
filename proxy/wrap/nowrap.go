package wrap

type NoWrapper struct {
}

func (h *NoWrapper) Wrap(data []byte) []byte {
	return data
}
