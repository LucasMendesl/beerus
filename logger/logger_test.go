package logger_test

import (
	"testing"

	"github.com/lucasmendesl/beerus/logger"
	"github.com/stretchr/testify/require"
)

type wantErr func(t *testing.T, err error) bool

func nopErr(t *testing.T, err error) bool {
	require.NoError(t, err)
	return false
}

func TestLoggerCreate(t *testing.T) {
	type args struct {
		logLevel string
        logFormat string
	}
	tests := []struct {
		name    string
		args    args
		wantErr wantErr
	}{
		{
			name: "invalid format",
			args: args{
                logFormat: "something",
			},
			wantErr: func(t *testing.T, err error) bool {
				require.EqualError(t, err, "invalid log formatter: something")
				return true
			},
		},
		{
			name: "invalid log level",
			args: args{
				logLevel: "something",
                logFormat: "text",
			},
			wantErr: func(t *testing.T, err error) bool {
				require.EqualError(t, err, "invalid log level: slog: level string \"something\": unknown name")
				return true
			},
		},
		{
			name: "create with text format",
			args: args{
                logLevel: "debug",
                logFormat: "text",
			},
			wantErr: nopErr,
		},
		{
			name: "create with json format",
			args: args{
                logLevel: "debug",
                logFormat: "json",
			},
			wantErr: nopErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := logger.Create(tt.args.logFormat, tt.args.logLevel)

			if tt.wantErr(t, err) {
				return
			}
		})
	}
}
