// +build !windows

package pdkshell

import (
	"reflect"
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
		name: "on non windows correct args are returned",
		args: args{
			pdkexe:  "/foo/lib/pdk",
			rubyexe: "/foo/lib/ruby",
		},
		want: []string{
			"/foo/lib/ruby", "-S", "--", "/foo/lib/pdk",
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
		name: "on *nix correct executable is returned",
		args: args{
			rubyexe: "/opt/puppetlabs/pdk/private/ruby/2.4.10/bin/ruby",
		},
		wantExecutable: "/opt/puppetlabs/pdk/private/ruby/2.4.10/bin/ruby",
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
