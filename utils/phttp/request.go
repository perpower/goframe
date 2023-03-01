package phttp

const (
	exceptionExit    = "exit"
	exceptionExitAll = "exit_all"
)

type Request struct{}

// Exit exits executing of current HTTP handler.
func Exit() {
	panic(exceptionExit)
}

// ExitAll exits executing of current and following HTTP handlers.
func ExitAll() {
	panic(exceptionExitAll)
}
