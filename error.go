package tcc

const (
	// ErrTryFailed means at least 1 error happened in Try phase,
	// but successfully canceled.
	ErrTryFailed = iota

	// ErrConfirmFailed means Try phase finished in success, but
	// at least 1 error happened in Confirm phase and never succeeded after some retries.
	// This should be never happened because if Try finished successfully,
	// Confirm must be finished successfully.
	// Basically, you need to fix inconsistent state manually
	ErrConfirmFailed

	// ErrCancelFailed means Try phase didn't finish in success,
	// and attempted to cancel all the services, but some resources could not be canceled.
	// Basically, you need to fix inconsistent state manually
	ErrCancelFailed
)

// Error knows what err happened in try/confirm/cancel phase.
type Error struct {
	failedPhase int
	err         error
	serviceName string
}

// FailedPhase returns FailedPhase code.
func (e *Error) FailedPhase() int {
	return e.failedPhase
}

// Error satisfies error interface
func (e *Error) Error() string {
	return e.err.Error()
}

// ServiceName returns the name of service which is failed to try/confirm/cancel.
func (e *Error) ServiceName() string {
	return e.serviceName
}
