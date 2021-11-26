package pct

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

// The top level structure
type RSS struct {
	Item []Item `xml:"channel>item" json:"news"`
}

// The structure of the information we wish to retrieve
type Item struct {
	Title string `xml:"title" json:"title"`
	Link  string `xml:"link" json:"link"`
}

func News(Url string, Format string) error {
	fmt.Printf("Retrieving HTML code of %s ...\n", Url)
	// Retrieve the html of the page
	resp, err := http.Get(Url) //nolint:gosec
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
	err = xml.Unmarshal([]byte(html), &items)
	if err != nil {
		return fmt.Errorf("Error ith retrieved yaml")
	}

	if Format == "table" {
		outputAsTable(items)
	} else if Format == "json" {
		_ = outputAsJson(items)
	}

	return nil
}

func outputAsTable(data RSS) {
	// Output the retrieved results as part of a table
	// Set the tables format
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	// Add the titles to the table
	tbl := table.New("Title", "Link")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	// Add each unmarshalled item as a new row
	for i := 0; i < len(data.Item); i++ {
		tbl.AddRow(data.Item[i].Title, data.Item[i].Link)
	}

	// Output the table
	tbl.Print()
}

func outputAsJson(data RSS) error {
	jsonData, _ := json.Marshal(data)

	var cleanedJsonData bytes.Buffer
	err := json.Indent(&cleanedJsonData, jsonData, "", "    ")
	if err != nil {
		return fmt.Errorf("Error formating Json")
	}
	fmt.Println(cleanedJsonData.String())

	return nil
}
