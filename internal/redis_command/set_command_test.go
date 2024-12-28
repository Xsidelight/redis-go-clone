package redis_command

import (
	"redis-go-clone/internal/model"
	"sync"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	tests := []struct {
		name     string
		cmdArray []any
		want     string
	}{
		{
			name:     "basic set",
			cmdArray: []any{"SET", "key", "value"},
			want:     "+OK\r\n",
		},
		{
			name:     "too few arguments",
			cmdArray: []any{"SET", "key"},
			want:     "-ERR wrong number of arguments for SET\r\n",
		},
		{
			name:     "invalid key type",
			cmdArray: []any{"SET", 123, "value"},
			want:     "-ERR invalid argument for SET\r\n",
		},
		{
			name:     "set with EX",
			cmdArray: []any{"SET", "key", "value", "EX", int64(60)},
			want:     "+OK\r\n",
		},
		{
			name:     "set with PX",
			cmdArray: []any{"SET", "key", "value", "PX", int64(1000)},
			want:     "+OK\r\n",
		},
		{
			name:     "set with EXAT",
			cmdArray: []any{"SET", "key", "value", "EXAT", time.Now().Add(time.Hour).Unix()},
			want:     "+OK\r\n",
		},
		{
			name:     "set with invalid expiry type",
			cmdArray: []any{"SET", "key", "value", "INVALID", int64(60)},
			want:     "-ERR unsupported expiry option\r\n",
		},
		{
			name:     "set with invalid expiry value",
			cmdArray: []any{"SET", "key", "value", "EX", "invalid"},
			want:     "-ERR invalid expiry value\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storedData := make(map[string]model.StoredData)
			mu := &sync.RWMutex{}
			if got := Set(tt.cmdArray, storedData, mu); got != tt.want {
				t.Errorf("Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseExpiryOptions(t *testing.T) {
	futureTime := time.Now().Add(time.Hour).Unix()

	tests := []struct {
		name        string
		options     []any
		wantErr     bool
		errContains string
	}{
		{
			name:        "too few options",
			options:     []any{"EX"},
			wantErr:     true,
			errContains: "invalid expiry arguments",
		},
		{
			name:        "invalid option type",
			options:     []any{123, int64(60)},
			wantErr:     true,
			errContains: "invalid expiry option type",
		},
		{
			name:        "invalid expiry value",
			options:     []any{"EX", "60"},
			wantErr:     true,
			errContains: "invalid expiry value",
		},
		{
			name:    "valid EX option",
			options: []any{"EX", int64(60)},
			wantErr: false,
		},
		{
			name:    "valid PX option",
			options: []any{"PX", int64(1000)},
			wantErr: false,
		},
		{
			name:    "valid EXAT option",
			options: []any{"EXAT", futureTime},
			wantErr: false,
		},
		{
			name:        "invalid EXAT value",
			options:     []any{"EXAT", int64(-1)},
			wantErr:     true,
			errContains: "invalid Unix time for EXAT",
		},
		{
			name:        "invalid PXAT value",
			options:     []any{"PXAT", int64(-1)},
			wantErr:     true,
			errContains: "invalid Unix time for PXAT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseExpiryOptions(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseExpiryOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errContains != "" && err.Error() != tt.errContains {
				t.Errorf("parseExpiryOptions() error = %v, want error containing %v", err, tt.errContains)
			}
		})
	}
}
