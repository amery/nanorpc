package nanorpc

func IsTimeout(err error) bool {
	if e, ok := err.(interface {
		IsTimeout() bool
	}); ok {
		return e.IsTimeout()
	}

	return false
}

func IsTemporary(err error) bool {
	if e, ok := err.(interface {
		IsTemporary() bool
	}); ok {
		return e.IsTemporary()
	}

	return false
}
