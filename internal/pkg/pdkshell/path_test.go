// +build !windows

package pdkshell

import "testing"

func Test_getPDKInstallDirectory(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "return correct install directory",
			want:    "/opt/puppetlabs/pdk",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPDKInstallDirectory(false)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPDKInstallDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getPDKInstallDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}
