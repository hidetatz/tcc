package tcc

import (
	"github.com/cenkalti/backoff"
	"golang.org/x/sync/errgroup"
)

// Option can set option to service
type Option func(s *Service)

// Service ...
type Service struct {
	try              func() error
	confirm          func() error
	cancel           func() error
	tried            bool
	trySucceeded     bool
	canceled         bool
	cancelSucceeded  bool
	confirmed        bool
	confirmSucceeded bool

	maxRetries uint64
}

// WithMaxRetries ...
func WithMaxRetries(maxRetries uint64) Option {
	return func(s *Service) {
		s.maxRetries = maxRetries
	}
}

// NewService returns service
func NewService(try, confirm, cancel func() error, opts ...Option) *Service {
	s := &Service{try: try, confirm: confirm, cancel: cancel}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Try ...
func (s *Service) Try() error { return s.try() }

// Confirm ...
func (s *Service) Confirm() error { return s.confirm() }

// Cancel ...
func (s *Service) Cancel() error { return s.cancel() }

// Orchestrator can orchestrate transactions
type Orchestrator struct {
	services []*Service
}

// NewOrchestrator returns interface Orchestrator
func NewOrchestrator(services []*Service, opts ...Option) *Orchestrator {
	for _, opt := range opts {
		for _, service := range services {
			opt(service)
		}
	}
	return &Orchestrator{services: services}
}

// Orchestrate can handle all the passed Service's transaction
func (o *Orchestrator) Orchestrate() error {
	if err := o.tryAll(); err != nil {
		return o.cancelAll()
	}
	return o.confirmAll()
}

func (o *Orchestrator) tryAll() error {
	eg := errgroup.Group{}
	backoff := backoff.NewExponentialBackOff()
	for _, s := range o.services {
		eg.Go(func() error {
			if s.maxRetries > 0 {
				_ := backoff.WithMaxRetries(backoff, s.maxRetries)
			}
			s.tried = true
			if err := s.Try(); err != nil {
				return err
			}
			s.trySucceeded = true
			return nil
		})
	}
	return eg.Wait()
}

func (o *Orchestrator) cancelAll() error {
	eg := errgroup.Group{}
	for _, s := range o.services {
		eg.Go(func() error {
			s.canceled = true
			if err := s.Cancel(); err != nil {
				return err
			}
			s.cancelSucceeded = true
			return nil
		})
	}
	return eg.Wait()
}

func (o *Orchestrator) confirmAll() error {
	eg := errgroup.Group{}
	for _, s := range o.services {
		eg.Go(func() error {
			s.confirmed = true
			if err := s.Confirm(); err != nil {
				return err
			}
			s.confirmSucceeded = true
			return nil
		})
	}
	return eg.Wait()
}
