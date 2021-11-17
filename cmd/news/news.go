package news

/*
Practise go by creating a function that retrieves an the information from the following BBC News xml link and gives a list of all the current headlines.
http://feeds.bbci.co.uk/news/technology/rss.xml
*/

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

const (
	default_url = "http://feeds.bbci.co.uk/news/technology/rss.xml"
)

var (
	Url string
)

// The top level structure
type RSS struct {
	Item []Item `xml:"channel>item"`
}

// The structure of the information we wish to retrieve
type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Guid        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
}

type NewsCommandI interface {
	CreateCommand() *cobra.Command
}

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:     "news",
		Short:   "Retrives the xml",
		Long:    `Retrives the xml from bbc news`,
		PreRunE: preExecute,
		RunE:    execute,
	}

	return tmp
}

func preExecute(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		Url = default_url
	}

	if len(args) == 1 {
		Url = args[0]
		_, err := url.ParseRequestURI(Url)
		if err != nil {
			return fmt.Errorf("The first argument should be a valid URL")
		}
	}

	if len(args) > 1 {
		return fmt.Errorf("Incorrect number of arguments; only a url can be passed")
	}

	return nil
}

func execute(cmd *cobra.Command, args []string) error {
	fmt.Printf("args ...\n", args)
	fmt.Printf("HTML code of %s ...\n", Url)
	// Retrieve the html of the page
	resp, err := http.Get(Url)
	if err != nil {
		return err
	}
	// do this now so it won't be forgotten
	defer resp.Body.Close()
	// reads html as a slice of bytes
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Output the raw html
	// fmt.Printf("%s\n", html)

	// Create a variable of the top structure type set out above
	var items RSS

	// Unmarshal the xml
	xml.Unmarshal([]byte(html), &items)

	// Output each item's Title and Link with a divider between them
	for i := 0; i < len(items.Item); i++ {
		fmt.Println("Title: " + items.Item[i].Title)
		fmt.Println("Link: " + items.Item[i].Link)
		fmt.Println("----------")
	}

	return nil
}
