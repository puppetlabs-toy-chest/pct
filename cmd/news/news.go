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

	"github.com/spf13/cobra"
)

const (
	url = "http://feeds.bbci.co.uk/news/technology/rss.xml"
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

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "news",
		Short: "Retrives the xml",
		Long:  `Retrives the xml from bbc news`,
		RunE:  execute,
	}

	return tmp
}

func execute(cmd *cobra.Command, args []string) error {
	fmt.Printf("HTML code of %s ...\n", url)
	// Retrieve the html of the page
	resp, err := http.Get(url)
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
