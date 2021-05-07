package utils

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestContains(t *testing.T) {
	type args struct {
		s   []string
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should return true item contained in list",
			args: args{
				s:   []string{"foo", "bar", "baz"},
				str: "foo",
			},
			want: true,
		},
		{
			name: "should return false if not contained in list",
			args: args{
				s:   []string{"foo", "bar", "baz"},
				str: "wakka",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.s, tt.args.str); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFind(t *testing.T) {
	type args struct {
		source []string
		match  string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "should find match",
			args: args{
				source: []string{"foo", "bar", "baz"},
				match:  "foo",
			},
			want: []string{"foo"},
		},
		{
			name: "should return nothing if no match",
			args: args{
				source: []string{"foo", "bar", "baz"},
				match:  "wakka",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Find(tt.args.source, tt.args.match); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetListOfFlags(t *testing.T) {
	type args struct {
		cmd           func() *cobra.Command
		argsV         []string
		flagsToIgnore []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "returns only the flag that was specified",
			want: []string{"--parallel"},
			args: args{
				cmd: func() *cobra.Command {
					c := &cobra.Command{
						Use:   "unit [flags]",
						Short: "Run unit tests",
						Long:  `Run unit tests`,
						RunE: func(cmd *cobra.Command, args []string) error {
							return nil
						},
					}
					c.Flags().Bool("parallel", false, "run unit tests in parallel")
					c.Flags().Bool("foo", false, "run unit tests in parallel")
					c.Flags().String("bar", "", "run unit tests in parallel")
					return c
				},
				argsV:         []string{"--parallel"},
				flagsToIgnore: []string{"log-level"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetListOfFlags(tt.args.cmd(), tt.args.argsV); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetListOfFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}

// :( All gone to waste

// type osFileOps interface {
// 	Getwd() (string, bool)
// 	Stat(string) (fs.FileInfo, error)
// }

// type osFO struct{}

// func(osFO) Getwd() (string, error){
// 	return "/foo/bar", nil
// }

// var statFilePathReturn fs.FileInfo


// func (osFO) Stat(string) (fs.FileInfo, error) {
// 	return statFilePathReturn, nil
// }

// func TestValidModuleRoot(t *testing.T) {
// 	var tests []struct {
// 		name 		string
// 		want 		string
// 		wantErr bool
// 	}

// 	tests = append(tests, struct{
// 		name 		string
// 		want 		string
// 		wantErr bool
// 	}{
// 			name: "Should return current working dir if valid module root",
// 			want: "/foo/bar",
// 			wantErr: false,
// 		},
// 	)
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			gotModuleRoot, err := ValidModuleRoot()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ValidModuleRoot() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if gotModuleRoot != tt.want {
// 				t.Errorf("ValidModuleRoot() = %v, want %v", gotModuleRoot, tt.want)
// 			}
// 		})
// 	}
// }
