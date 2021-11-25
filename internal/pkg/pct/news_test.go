package pct_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type NewsTest struct {
	name           string
	args           newsArgs
	expected       newsExpected
}

// what goes in
type newsArgs struct {
	url 	string
	format	string
}

// what comes out
type newsExpected struct {
	errorMsg 	string
	fmt_println string
}

func TestNews(t *testing.T) {

	tests := []NewsTest{
		{
			name: "valid url and table flag",
			args: newsArgs{
				url: "http://feeds.bbci.co.uk/news/technology/rss.xml",
				format: "table",
			},
			expected: newsExpected{
				errorMsg: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.expected.errorMsg != "" && err != nil {
				assert.Contains(t, err.Error(), tt.expected.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
