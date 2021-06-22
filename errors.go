package bs

var fnErrorHandler = func(err error) {
	panic(err)
}

// SetErrorHandler sets the behavior when an error is encountered while running most commands.
// The default behavior is to panic.
func SetErrorHandler(fnErr func(err error)) {
	Verbose("Error handler changed")
	fnErrorHandler = fnErr
}
