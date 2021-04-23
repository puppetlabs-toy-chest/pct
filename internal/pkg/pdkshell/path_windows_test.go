// +build windows

package pdkshell

import (
	"regexp"
	"testing"
)

// we can't do this yet because we need to mock getRegistryStringKey
// func Test_getPDKInstallDirectory(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		want    string
// 		wantErr bool
// 	}{
// 		{
// 			name:    "return correct install directory",
// 			want:    "C:\\Program Files\\Puppet Labs\\DevelopmentKit\\",
// 			wantErr: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := getPDKInstallDirectory(false)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("getPDKInstallDirectory() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("getPDKInstallDirectory() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func Test_getRegistryStringKey(t *testing.T) {
	type args struct {
		path string
		key  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "finds a generic registry entry",
			args: args{
				path: `SOFTWARE\Microsoft\Windows NT\CurrentVersion`,
				key:  "SystemRoot",
			},
			want:    `(?i)C\:\\windows`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRegistryStringKey(tt.args.path, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRegistryStringKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			r := regexp.MustCompile(tt.want)
			if !r.MatchString(got) {
				t.Errorf("output did not match. regexp: /%s/ output: %s", r, got)
				return
			}
		})
	}
}

func Test_getShortPath(t *testing.T) {
	type args struct {
		longPath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "basic longpath converted to short path",
			args: args{
				longPath: `C:\Program Files\WindowsPowerShell`,
			},
			want: `C:\PROGRA~1\WindowsPowerShell`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getShortPath(tt.args.longPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("getShortPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getShortPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
