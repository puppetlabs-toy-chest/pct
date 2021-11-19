package news

/*
Practise go by creating a function that retrieves an the information from the following BBC News xml link and gives a list of all the current headlines.
http://feeds.bbci.co.uk/news/technology/rss.xml
*/

import (
	"fmt"
	"net/url"

	"github.com/puppetlabs/pdkgo/internal/pkg/pct"
	"github.com/spf13/cobra"
)

const (
	default_url = "http://feeds.bbci.co.uk/news/technology/rss.xml"
)

var (
	Url          string
	outputFormat string
)

type NewsCommandI interface {
	CreateCommand() *cobra.Command
}

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:     "news <url> [flags]",
		Short:   "Retrives the xml",
		Long:    `Retrives the xml from bbc news`,
		PreRunE: preExecute,
		RunE:    execute,
	}

	tmp.Flags().SortFlags = false
	tmp.Flags().StringVarP(&outputFormat, "format", "f", "table", "the desired format of the output; either json or table.")

	return tmp
}

// Verify that the input given is as exected
func preExecute(cmd *cobra.Command, args []string) error {
	// If no url is given, default to the stored on
	if len(args) < 1 {
		Url = default_url
	}

	// If a url is given valify that it is valid
	if len(args) == 1 {
		Url = args[0]
		_, err := url.ParseRequestURI(Url)
		if err != nil {
			return fmt.Errorf("The first argument should be a valid URL")
		}
	}

	// If more than one input is given then throw an error
	if len(args) > 1 {
		return fmt.Errorf("Incorrect number of arguments; only a url can be passed")
	}

	return nil
}

func execute(cmd *cobra.Command, args []string) error {
	pct.News(Url, outputFormat)
	return nil
}
