package pct_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/puppetlabs/pdkgo/internal/pkg/mock"
	"github.com/puppetlabs/pdkgo/internal/pkg/pct"
	"github.com/stretchr/testify/assert"
)

const html = `<?xml version="1.0" encoding="UTF-8"?>
<?xml-stylesheet title="XSL_formatting" type="text/xsl" href="/shared/bsp/xsl/rss/nolsol.xsl"?>
<rss xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:content="http://purl.org/rss/1.0/modules/content/" xmlns:atom="http://www.w3.org/2005/Atom" version="2.0" xmlns:media="http://search.yahoo.com/mrss/">
	<channel>
		<title><![CDATA[BBC News - Technology]]></title>
		<description><![CDATA[BBC News - Technology]]></description>
		<link>https://www.bbc.co.uk/news/</link>
		<image>
			<url>https://news.bbcimg.co.uk/nol/shared/img/bbc_news_120x60.gif</url>
			<title>BBC News - Technology</title>
			<link>https://www.bbc.co.uk/news/</link>
		</image>
		<generator>RSS for Node</generator>
		<lastBuildDate>Thu, 02 Dec 2021 10:58:36 GMT</lastBuildDate>
		<copyright><![CDATA[Copyright: (C) British Broadcasting Corporation, see http://news.bbc.co.uk/2/hi/help/rss/4498287.stm for terms and conditions of reuse.]]></copyright>
		<language><![CDATA[en-gb]]></language>
		<ttl>15</ttl>
		<item>
			<title><![CDATA[Cryptocurrency executives to be questioned in Congress]]></title>
			<description><![CDATA[Washington is joining other governments in scrutinising the rapidly expanding sector more closely.]]></description>
			<link>https://www.bbc.co.uk/news/business-59496509?at_medium=RSS&amp;at_campaign=KARANGA</link>
			<guid isPermaLink="false">https://www.bbc.co.uk/news/business-59496509</guid>
			<pubDate>Thu, 02 Dec 2021 00:02:25 GMT</pubDate>
		</item>
	</channel>
</rss>`

type NewsTest struct {
	name     string
	args     newsArgs
	expected newsExpected
}

// what goes in
type newsArgs struct {
	url    string
	format string
}

// what comes out
type newsExpected struct {
	httpResponse   int
	httpOutput     string
	errorExpected  bool
	expectedOutput string
}

func TestNews(t *testing.T) {

	tests := []NewsTest{
		{
			name: "valid html and table flag",
			args: newsArgs{
				url:    "http://feeds.bbci.co.uk/news/technology/rss.xml",
				format: "table",
			},
			expected: newsExpected{
				httpResponse:  200,
				httpOutput:    html,
				errorExpected: false,
			},
		},
		{
			name: "valid html and json flag",
			args: newsArgs{
				url:    "http://feeds.bbci.co.uk/news/technology/rss.xml",
				format: "json",
			},
			expected: newsExpected{
				httpResponse:  200,
				httpOutput:    html,
				errorExpected: false,
			},
		},
		{
			name: "failed http.get",
			args: newsArgs{
				url:    "http://feeds.bbci.co.uk/news/technology/r",
				format: "json",
			},
			expected: newsExpected{
				httpResponse:   400,
				httpOutput:     "@quality",
				errorExpected:  true,
				expectedOutput: "Web request error",
			},
		},
		{
			name: "invalid html returned",
			args: newsArgs{
				url:    "http://feeds.bbci.co.uk/news/technology/rss.xml",
				format: "json",
			},
			expected: newsExpected{
				httpResponse:   200,
				httpOutput:     "@quality",
				errorExpected:  true,
				expectedOutput: "Web request error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up a mock for the http get
			nc := pct.NewsCommand{
				HttpClient: &mock.HTTPClient{
					RequestResponse: &http.Response{
						StatusCode: tt.expected.httpResponse,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(tt.expected.httpOutput))),
					},
					ErrResponse: tt.expected.errorExpected,
				},
			}

			// Once the above mock is done the rest should run as expected
			//	Check to ensure the correct method is called next, how can I do this?
			err := nc.News(tt.args.url, tt.args.format)
			if tt.expected.errorExpected == true {
				// Check that it errored as expected
				assert.Error(t, err)
				// Check that the returned error was correct
				assert.Equal(t, tt.expected.expectedOutput, err.Error())
				return
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
