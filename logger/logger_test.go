package logger_test

import (
	"testing"

	"github.com/lucasmendesl/beerus/config"
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
        config config.Logging
	}
	tests := []struct {
		name    string
		args    args
		wantErr wantErr
	}{
		{
			name: "invalid format",
			args: args{
                config: config.Logging{
                    Format: "something",
                },
			},
			wantErr: func(t *testing.T, err error) bool {
				require.EqualError(t, err, "invalid log formatter: something")
				return true
			},
		},
		{
			name: "invalid log level",
			args: args{
                config: config.Logging{
                    Level: "something",
                    Format: "text",
                },
			},
			wantErr: func(t *testing.T, err error) bool {
				require.EqualError(t, err, "invalid log level: slog: level string \"something\": unknown name")
				return true
			},
		},
		{
			name: "create with text format",
			args: args{
                config: config.Logging{
                    Format: "text",
                    Level: "debug",
                },
			},
			wantErr: nopErr,
		},
		{
			name: "create with json format",
			args: args{
                config: config.Logging{
                    Format: "json",
                    Level: "debug",
                },
			},
			wantErr: nopErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := logger.Create(tt.args.config)

			if tt.wantErr(t, err) {
				return
			}
		})
	}
}
