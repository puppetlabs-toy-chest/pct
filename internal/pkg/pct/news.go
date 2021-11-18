package pct

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fatih/color"
	"github.com/rodaine/table"
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

func News(Url string) error {
	fmt.Printf("Retrieving HTML code of %s ...\n", Url)
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

	// Create a variable of the top structure type set out above
	var items RSS

	// Unmarshal the xml
	xml.Unmarshal([]byte(html), &items)

	// Output the retrieved results as part of a table
	// Set the tables format
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	// Add the titles to the table
	tbl := table.New("Title", "Link")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	// Add each unmarshalled item as a new row
	for i := 0; i < len(items.Item); i++ {
		tbl.AddRow(items.Item[i].Title, items.Item[i].Link)
	}

	// Output the table
	tbl.Print()

	return nil
}
