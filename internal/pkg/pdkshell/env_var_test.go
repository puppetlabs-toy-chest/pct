package pdkshell

import (
	"reflect"
	"testing"
)

type osHelpersImplMock struct{}

var mockedEnvironReturn []string

func (osHelpersImplMock) Environ() []string {
	return mockedEnvironReturn
}

func Test_getEnvironVars(t *testing.T) {

	osUtil = osHelpersImplMock{}

	tests := []struct {
		name                string
		mockedEnvironReturn []string
		want                []string
	}{
		{
			name:                "Drops ENV VAR to unset and passes through permitted",
			mockedEnvironReturn: []string{"GEM_HOME=/foo/bar", "VALID_VAR=/baz/qux"},
			want:                []string{"VALID_VAR=/baz/qux"},
		},
		{
			name:                "Handles multiple instances of the same ENV VAR to be dropped",
			mockedEnvironReturn: []string{"RUBYPATH=/foo/bar", "VALID_VAR=/baz/qux", "RUBYPATH=/foo/bar/bizz/buzz"},
			want:                []string{"VALID_VAR=/baz/qux"},
		},
	}

	for _, tt := range tests {
		mockedEnvironReturn = tt.mockedEnvironReturn
		t.Run(tt.name, func(t *testing.T) {
			if got := getEnvironVars(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getEnvironVars() = %v, want %v", got, tt.want)
			}
		})
	}
}
