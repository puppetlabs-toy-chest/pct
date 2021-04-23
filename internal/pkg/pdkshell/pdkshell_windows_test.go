// +build windows

package pdkshell

import (
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func Test_buildCommandArgs(t *testing.T) {
	type args struct {
		args    []string
		rubyexe string
		pdkexe  string
	}

	var tests []struct {
		name string
		args args
		want []string
	}

	tests = append(tests, struct {
		name string
		args args
		want []string
	}{
		name: "windows test",
		args: args{
			args:    []string{},
			pdkexe:  "c:/foo/lib/pdk",
			rubyexe: "c:/foo/lib/ruby",
		},
		want: []string{
			"/c", "c:/foo/lib/ruby", "-S", "--", "c:/foo/lib/pdk",
		},
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildCommandArgs(tt.args.args, tt.args.rubyexe, tt.args.pdkexe); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildCommandArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildExecutable(t *testing.T) {
	type args struct {
		rubyexe string
	}
	var tests []struct {
		name           string
		args           args
		wantExecutable string
	}
	tests = append(tests, struct {
		name           string
		args           args
		wantExecutable string
	}{
		name: "on windows correct executable is returned",
		args: args{
			rubyexe: runtime.GOOS,
		},
		// wantExecutable: "C:\\Windows\\system32\\cmd.exe",
		wantExecutable: `cmd.exe`,
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExecutable := buildExecutable(tt.args.rubyexe)
			if !strings.HasSuffix(gotExecutable, tt.wantExecutable) {
				t.Errorf("buildExecutable() = %v, want %v", gotExecutable, tt.wantExecutable)
			}
		})
	}
}
