package tcc

import (
	"errors"
	"sync"

	"github.com/cenkalti/backoff"
	"golang.org/x/sync/errgroup"
)

// Option can set option to service
// Option can be passed to NewService() and NewOrchestrator,
// if you pass it to both, the one which is passed to NewOrchestrator will be used
type Option func(s *orchestrator)

// WithMaxRetries sets limitation of retry times
func WithMaxRetries(maxRetries uint64) Option {
	return func(o *orchestrator) {
		o.backoff = backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxRetries)
	}
}

// Orchestrator can orchestrate multiple service
// First, call every service's try() asynchronously.
// If all the try succeeded, call every service's confirm().
// If even one of the services' try fails, every service's cancel will be called.
type Orchestrator interface {
	Orchestrate() error
}

type orchestrator struct {
	services []*Service
	backoff  backoff.BackOff

	sync.Mutex
}

// NewOrchestrator returns interface Orchestrator
func NewOrchestrator(services []*Service, opts ...Option) Orchestrator {
	maxRetries := uint64(10)
	o := &orchestrator{
		services: services,
		backoff:  backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxRetries),
		Mutex:    sync.Mutex{},
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Orchestrate can handle all the passed Service's transaction
func (o *orchestrator) Orchestrate() error {
	if tryErr := o.tryAll(); tryErr != nil {
		if cancelErr := o.cancelAll(); cancelErr != nil {
			return cancelErr
		}
		return tryErr
	}
	return o.confirmAll()
}

func (o *orchestrator) tryAll() error {
	eg := errgroup.Group{}
	for _, s := range o.services {
		s := s
		eg.Go(func() error {
			s.tried = true
			err := s.Try()
			if err != nil {
				return &Error{
					failedPhase: ErrTryFailed,
					err:         err,
					serviceName: s.name,
				}
			}
			s.trySucceeded = true
			return nil
		})
	}
	return eg.Wait()
}

func (o *orchestrator) confirmAll() error {
	eg := errgroup.Group{}
	for _, s := range o.services {
		s := s
		eg.Go(func() error {
			s.confirmed = true
			if !s.trySucceeded {
				return &Error{
					failedPhase: ErrConfirmFailed,
					err:         errors.New("try did not succeed"),
					serviceName: s.name,
				}
			}
			o.Lock()
			defer o.Unlock()
			err := backoff.Retry(s.Confirm, o.backoff)
			if err != nil {
				return &Error{
					failedPhase: ErrConfirmFailed,
					err:         err,
					serviceName: s.name,
				}
			}
			s.confirmSucceeded = true
			return nil
		})
	}
	return eg.Wait()
}

func (o *orchestrator) cancelAll() error {
	eg := errgroup.Group{}
	for _, s := range o.services {
		s := s
		eg.Go(func() error {
			if !s.tried {
				return nil
			}
			s.canceled = true
			o.Lock()
			defer o.Unlock()
			err := backoff.Retry(s.Cancel, o.backoff)
			if err != nil {
				return &Error{
					failedPhase: ErrCancelFailed,
					err:         err,
					serviceName: s.name,
				}
			}
			s.cancelSucceeded = true
			return nil
		})
	}
	return eg.Wait()
}
