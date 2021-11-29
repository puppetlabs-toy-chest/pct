package news

/*
Practise go by creating a function that retrieves an the information from the following BBC News xml link and gives a list of all the current headlines.
http://feeds.bbci.co.uk/news/technology/rss.xml
*/

import (
	"fmt"

	"github.com/puppetlabs/pdkgo/internal/pkg/pct"
	"github.com/spf13/cobra"
)

var (
	Url          string
	OutputFormat string

	Urls = map[string]string{
		"bbc":        "http://feeds.bbci.co.uk/news/technology/rss.xml",
		"nasa":       "https://www.nasa.gov/rss/dyn/breaking_news.rss",
		"nytimes":    "https://rss.nytimes.com/services/xml/rss/nyt/US.xml",
		"lifehacker": "https://lifehacker.com/rss",
	}
)

type NewsCommandI interface {
	CreateCommand() *cobra.Command
}

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:       "news <url> [flags]",
		Short:     "Retrives the xml from selet websites",
		Long:      `Retrives the xml from a selection of websites, valid options are 'bbc', 'nasa', 'nytimes' and 'lifehacker'`,
		PreRunE:   preExecute,
		RunE:      execute,
		ValidArgs: []string{"bbc", "nasa", "nytimes", "lifehacker"},
	}

	tmp.Flags().SortFlags = false
	tmp.Flags().StringVarP(&OutputFormat, "format", "f", "table", "the desired format of the output; either json or table.")

	return tmp
}

// Verify that the input given is as exected
func preExecute(cmd *cobra.Command, args []string) error {
	// If no url is given, default to the stored on
	if len(args) < 1 {
		Url = Urls["bbc"]
	}

	// If a url is given valify that it is valid
	if len(args) == 1 {
		Url = Urls[args[0]]
		if Url == "" {
			return fmt.Errorf("The first argument should be a valid URL key")
		}
	}

	// If more than one input is given then throw an error
	if len(args) > 1 {
		return fmt.Errorf("Incorrect number of arguments; only a url can be passed")
	}

	// If flag is invalid
	if OutputFormat != "table" && OutputFormat != "json" {
		return fmt.Errorf("Invalid value for format flag has been given")
	}

	return nil
}

func execute(cmd *cobra.Command, _ []string) error {
	_ = pct.News(Url, OutputFormat)
	return nil
}
