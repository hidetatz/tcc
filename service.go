package tcc

// Service can be TCC service, which can Try(), Confirm(), and Cancel()
type Service struct {
	name string

	try     func() error
	confirm func() error
	cancel  func() error

	tried            bool
	trySucceeded     bool
	confirmed        bool
	confirmSucceeded bool
	canceled         bool
	cancelSucceeded  bool
}

// NewService returns service with passed functions
func NewService(name string, try, confirm, cancel func() error) *Service {
	return &Service{name: name, try: try, confirm: confirm, cancel: cancel}
}

// Try executes passed try function.
// In try phase, service will do some reservation or precondition satisfyment.
// After try phase is finished successfully, Confirm called.
// Try can fail, but if try succeeded, confirm must succeed.
// If try fails, Cancel will be called.
// Try never be retried.
func (s *Service) Try() error { return s.try() }

// Confirm executes passed confirm function.
// In confirm phase, service will confirm things which is reserved in try phase.
// Basically Confirm should never return error, except network or infrastructure issues.
// This will be retried 10 times by default.
func (s *Service) Confirm() error { return s.confirm() }

// Cancel executes passed cancel function.
// This will be called after Try phase failed.
// In Cancel phase, service will revert the state which is changed by try phase.
// Basically Confirm should never return error, except network or infrastructure issues.
// This will be retried 10 times by default.
func (s *Service) Cancel() error { return s.cancel() }

// Tried returns if the service try() called
func (s *Service) Tried() bool {
	return s.tried
}

// TrySucceeded returns if the service try() succeeded
func (s *Service) TrySucceeded() bool {
	return s.trySucceeded
}

// Confirmed returns if the service confirm() called
func (s *Service) Confirmed() bool {
	return s.tried
}

// ConfirmtSucceeded returns if the service confirm() succeeded
func (s *Service) ConfirmSucceeded() bool {
	return s.tried
}

// Canceled returns if the service cancel() is called
func (s *Service) Canceled() bool {
	return s.tried
}

// CancelSucceeded returns if the service cancel() succeeded
func (s *Service) CancelSucceeded() bool {
	return s.tried
}
