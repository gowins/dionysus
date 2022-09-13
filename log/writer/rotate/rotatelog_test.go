package rotate

import (
	"fmt"
	"path"
	"testing"
	"time"
)

func TestNewRotateLogs(t *testing.T) {
	testDefaultLogDir := "/tmp"
	type args struct {
		root     string
		sub      string
		fn       string
		age      time.Duration
		duration time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Normal",
			args: args{
				root:     testDefaultLogDir,
				sub:      "debug",
				fn:       "test.log",
				age:      DefaultRotateMaxAge,
				duration: DefaultRotateDuration,
			},
			wantErr: false,
		},
		{
			name: "Duration Minute",
			args: args{
				root:     testDefaultLogDir,
				sub:      "debug",
				fn:       "test.log",
				age:      DefaultRotateMaxAge,
				duration: time.Minute,
			},
			wantErr: false,
		},
		{
			name: "Duration Second",
			args: args{
				root:     testDefaultLogDir,
				sub:      "debug",
				fn:       "test.log",
				age:      DefaultRotateMaxAge,
				duration: time.Second,
			},
			wantErr: false,
		},
		//{
		//	name: "Can't Create Dir",
		//	args: args{
		//		root:     "/bin",
		//		sub:      "debug",
		//		fn:       "test.log",
		//		age:      DefaultRotateMaxAge,
		//		duration: time.Second,
		//	},
		//	wantErr: true,
		//},
	}
	for _, c := range tests {
		tt := c
		t.Run(tt.name, func(t *testing.T) {
			cfg := NewWriterConfig()
			cfg.Dir = tt.args.root
			cfg.Sub = tt.args.sub
			cfg.Filename = tt.args.fn
			cfg.Age = tt.args.age
			cfg.Duration = tt.args.duration
			got, err := NewRotateLogger(cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRotateLogs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil {
				n, err := got.Write([]byte("123"))
				if err != nil {
					t.Error(err)
				}
				if n != 3 {
					t.Error("n != 3")
				}
				fn := got.CurrentFileName()
				now := time.Now().In(Location())

				var suffix string
				if tt.args.duration < time.Hour {
					suffix = fmt.Sprintf(".%d%02d%02d%02d%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute())
				} else if tt.args.duration < DefaultRotateDuration {
					suffix = fmt.Sprintf(".%d%02d%02d%02d00", now.Year(), now.Month(), now.Day(), now.Hour())
				} else {
					suffix = fmt.Sprintf(".%d%02d%02d0000", now.Year(), now.Month(), now.Day())
				}

				actual := path.Join(tt.args.root, tt.args.sub, tt.args.fn+suffix)
				if fn != actual {
					t.Errorf("fn != actual,fn: %s,actual: %s", fn, actual)
				}
			}
		})
	}
}
