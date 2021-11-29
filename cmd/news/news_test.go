package news_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/puppetlabs/pdkgo/cmd/news"
	"github.com/stretchr/testify/assert"
)

func TestCreateCommand(t *testing.T) {
	tests := []struct {
		name    	string
		args    	[]string
		out     	string
		expectErr 	bool
	}{
		{
			name:    "executes without error",
			out:     "",
			expectErr: false,
		},
		{
			name:    "executes without error for valid url key",
			args:    []string{"bbc"},
			out:     "",
			expectErr: false,
		},
		{
			name:    "executes without error for valid url key and format flag table",
			args:    []string{"bbc", "-f", "table"},
			out:     "",
			expectErr: false,
		},
		{
			name:    "executes without error for valid url key and format flag json",
			args:    []string{"bbc", "-f", "json"},
			out:     "",
			expectErr: false,
		},
		{
			name:    "executes with error for invalid url key",
			args:    []string{"bbs"},
			out:     "Error: The first argument should be a valid URL key\n",
			expectErr: true,
		},
		{
			name:    "executes with error for to many arguments",
			args:    []string{"bbc", "nasa"},
			out:     "Error: Incorrect number of arguments; only a url can be passed\n",
			expectErr: true,
		},
		{
			name:    "executes with error for invalid format flag",
			args:    []string{"-f", "jso"},
			out:     "Error: Invalid value for format flag has been given\n",
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := news.CreateCommand()
			b := bytes.NewBufferString("")
			cmd.SetErr(b)
			cmd.SetArgs(tt.args)

			// If an error is expected check that it is correct
			err := cmd.Execute()

			if tt.expectErr {
				assert.Error(t, err)
				out, _ := ioutil.ReadAll(b)
				assert.Equal(t, string(out), tt.out)
			} else {
				assert.NoError(t, err)
				out, _ := ioutil.ReadAll(b)
				assert.Equal(t, string(out), tt.out)
			}
		})
	}
}
