package logger

import "testing"

func TestNewLogger_DefaultFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "keep configured filename",
			filename: "/tmp/custom.log",
			want:     "/tmp/custom.log",
		},
		{
			name:     "use app default when empty",
			filename: "",
			want:     "./logs/app.log",
		},
		{
			name: "use app default for all services when empty filename",
			want: "./logs/app.log",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Filename: tt.filename,
				Env:      "test",
			}

			_ = NewLogger(cfg)
			got := cfg.Filename
			if got != tt.want {
				t.Fatalf("NewLogger() filename = %q, want %q", got, tt.want)
			}
		})
	}
}
