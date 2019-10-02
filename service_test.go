package tcc

import (
	"testing"
)

func TestService_Try(t *testing.T) {
	var tried bool
	type fields struct {
		try func() error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "try",
			fields: fields{
				try: func() error {
					tried = true
					return nil
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				try: tt.fields.try,
			}
			if err := s.Try(); (err != nil) != tt.wantErr {
				t.Errorf("Service.Try() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tried {
				t.Errorf("try() is not called")
			}
		})
	}
}

func TestService_Confirm(t *testing.T) {
	var confirmed bool
	type fields struct {
		confirm func() error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "confirm",
			fields: fields{
				confirm: func() error {
					confirmed = true
					return nil
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				confirm: tt.fields.confirm,
			}
			if err := s.Confirm(); (err != nil) != tt.wantErr {
				t.Errorf("Service.Confirm() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !confirmed {
				t.Errorf("confirm() is not called")
			}
		})
	}
}

func TestService_Cancel(t *testing.T) {
	var canceled bool
	type fields struct {
		cancel func() error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "confirm",
			fields: fields{
				cancel: func() error {
					canceled = true
					return nil
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				cancel: tt.fields.cancel,
			}
			if err := s.Cancel(); (err != nil) != tt.wantErr {
				t.Errorf("Service.Cancel() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !canceled {
				t.Errorf("cancel() is not called")
			}
		})
	}
}

func TestService_Params(t *testing.T) {
	type fields struct {
		name             string
		tried            bool
		trySucceeded     bool
		confirmed        bool
		confirmSucceeded bool
		canceled         bool
		cancelSucceeded  bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "true",
			fields: fields{
				tried:            true,
				trySucceeded:     true,
				confirmed:        true,
				confirmSucceeded: true,
				canceled:         true,
				cancelSucceeded:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				tried:            tt.fields.tried,
				trySucceeded:     tt.fields.trySucceeded,
				confirmed:        tt.fields.confirmed,
				confirmSucceeded: tt.fields.confirmSucceeded,
				canceled:         tt.fields.canceled,
				cancelSucceeded:  tt.fields.cancelSucceeded,
			}
			if got := s.Tried(); got != tt.fields.tried {
				t.Errorf("Service.Tried() = %v, want %v", got, tt.fields.tried)
			}
			if got := s.TrySucceeded(); got != tt.fields.trySucceeded {
				t.Errorf("Service.TrySucceeded() = %v, want %v", got, tt.fields.trySucceeded)
			}
			if got := s.Confirmed(); got != tt.fields.confirmed {
				t.Errorf("Service.Confirmed() = %v, want %v", got, tt.fields.confirmed)
			}
			if got := s.ConfirmSucceeded(); got != tt.fields.confirmSucceeded {
				t.Errorf("Service.ConfirmSucceeded() = %v, want %v", got, tt.fields.confirmSucceeded)
			}
			if got := s.Canceled(); got != tt.fields.canceled {
				t.Errorf("Service.Canceled() = %v, want %v", got, tt.fields.canceled)
			}
			if got := s.CancelSucceeded(); got != tt.fields.cancelSucceeded {
				t.Errorf("Service.CancelSucceeded() = %v, want %v", got, tt.fields.cancelSucceeded)
			}
		})
	}
}
