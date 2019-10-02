package tcc

import (
	"errors"
	"testing"
)

func TestError_FailedPhase(t *testing.T) {
	type fields struct {
		failedPhase int
		err         error
		serviceName string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "Try failed",
			fields: fields{failedPhase: ErrTryFailed},
			want:   ErrTryFailed,
		},
		{
			name:   "Confirm failed",
			fields: fields{failedPhase: ErrConfirmFailed},
			want:   ErrConfirmFailed,
		},
		{
			name:   "Cancel failed",
			fields: fields{failedPhase: ErrCancelFailed},
			want:   ErrCancelFailed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				failedPhase: tt.fields.failedPhase,
				err:         tt.fields.err,
				serviceName: tt.fields.serviceName,
			}
			if got := e.FailedPhase(); got != tt.want {
				t.Errorf("Error.FailedPhase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	type fields struct {
		failedPhase int
		err         error
		serviceName string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "error happened",
			fields: fields{err: errors.New("test")},
			want:   "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				failedPhase: tt.fields.failedPhase,
				err:         tt.fields.err,
				serviceName: tt.fields.serviceName,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_ServiceName(t *testing.T) {
	type fields struct {
		failedPhase int
		err         error
		serviceName string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "error happened",
			fields: fields{serviceName: "test"},
			want:   "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				failedPhase: tt.fields.failedPhase,
				err:         tt.fields.err,
				serviceName: tt.fields.serviceName,
			}
			if got := e.ServiceName(); got != tt.want {
				t.Errorf("Error.ServiceName() = %v, want %v", got, tt.want)
			}
		})
	}
}
