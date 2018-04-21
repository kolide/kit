package actor

// Actor is a struct that can be used to represent interruptible workloads
type Actor struct {
	Execute   func() error
	Interrupt func(error)
}
