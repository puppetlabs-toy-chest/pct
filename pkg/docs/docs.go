package docs

import (
	"embed"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/gernest/front"
	"github.com/puppetlabs/pct/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
)

type Docs struct {
	AFS             *afero.Afero
	IOFS            *afero.IOFS
	TermRenderer    *glamour.TermRenderer
	DocsFileSystem  *embed.FS
	MarkdownHandler *front.Matter
	ParsedDocsCache []MarkdownDoc
}

type MarkdownDoc struct {
	FrontMatter DocsFrontMatter
	Body        string
}

type Title struct {
	Short string
	Long  string
}

type DocsFrontMatter struct {
	Title       Title
	Description string
	Category    string
	Tags        []string
}

type DocsI interface {
	ReadMarkdownDoc(text string) (mdc MarkdownDoc)
	InitRenderer() (err error)
	Render(body string) (output string, err error)
	ListByCategory(category string)
	ListByTag(tag string)
	List()
	FindAndParse(docsFolderPath string)
}

func (d *Docs) InitRenderer() (err error) {
	if d.TermRenderer == nil {
		d.TermRenderer, err = glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(120),
		)
		return err
	}
	return err
}

func (d *Docs) Render(body string) (output string, err error) {
	err = d.InitRenderer()
	if err != nil {
		return "", err
	}
	output, err = d.TermRenderer.Render(body)
	return output, err
}

func (d *Docs) InitHandler() {
	if d.MarkdownHandler == nil {
		d.MarkdownHandler = front.NewMatter()
		d.MarkdownHandler.Handle("---", front.YAMLHandler)
	}
}

func (d *Docs) FindAndParse(docsFolderPath string) {
	// ignore errors for now
	dirEntries, _ := d.DocsFileSystem.ReadDir(docsFolderPath)
	for _, entry := range dirEntries {
		entryPath := fmt.Sprintf("%s/%s", docsFolderPath, entry.Name())

		if entry.IsDir() {
			d.FindAndParse(entryPath)
		} else {
			log.Debug().Msgf("Parsing Documentation File: %s", entryPath)
			// read file, ignoring errors for now
			raw, _ := d.DocsFileSystem.ReadFile(entryPath)
			// parse and append to cache
			d.InitHandler()
			fm, b, err := d.MarkdownHandler.Parse(strings.NewReader(string(raw)))
			if err != nil {
				log.Warn().Msgf("Could not parse %s", entryPath)
				// Some docs might need to be skipped in terminal
			} else if skip, _ := fm["skipTerminal"]; skip != true { //nolint
				// Turn the tags into an array of strings for further use
				tagsAsInterface := fm["tags"].([]interface{})
				tags := make([]string, len(tagsAsInterface))
				for i, v := range tagsAsInterface {
					tags[i] = v.(string)
				}
				d.ParsedDocsCache = append(d.ParsedDocsCache, MarkdownDoc{
					Body: b,
					FrontMatter: DocsFrontMatter{
						Title: Title{
							Short: strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name())),
							Long:  fm["title"].(string),
						},
						Description: fm["description"].(string),
						Category:    fm["category"].(string),
						Tags:        tags,
					},
				})
			}
		}
	}
}

func (d *Docs) ListTags(docs []MarkdownDoc) (tags []string) {
	for _, doc := range docs {
		for _, tag := range doc.FrontMatter.Tags {
			if !utils.Contains(tags, tag) {
				tags = append(tags, tag)
			}
		}
	}
	return tags
}

func (d *Docs) ListCategories(docs []MarkdownDoc) (categories []string) {
	for _, doc := range docs {
		if !utils.Contains(categories, doc.FrontMatter.Category) {
			categories = append(categories, doc.FrontMatter.Category)
		}
	}
	return categories
}

func (d *Docs) ListTitles(docs []MarkdownDoc) (titles []Title) {
	for _, doc := range docs {
		titles = append(titles, doc.FrontMatter.Title)
	}
	return titles
}

func (d *Docs) CompleteTitle(docs []MarkdownDoc, match string) []string {
	var titles []string
	for _, title := range d.ListTitles(docs) {
		if strings.HasPrefix(title.Short, match) {
			titles = append(titles, fmt.Sprintf("%s\t%s", title.Short, title.Long))
		}
	}
	return titles
}

func (d *Docs) FilterByTag(tag string, docs []MarkdownDoc) (filteredDocs []MarkdownDoc) {
	for _, doc := range docs {
		if utils.Contains(doc.FrontMatter.Tags, tag) {
			filteredDocs = append(filteredDocs, doc)
		}
	}
	return filteredDocs
}

func (d *Docs) FilterByCategory(category string, docs []MarkdownDoc) (filteredDocs []MarkdownDoc) {
	for _, doc := range docs {
		if doc.FrontMatter.Category == category {
			filteredDocs = append(filteredDocs, doc)
		}
	}
	return filteredDocs
}

func (d *Docs) SelectDocument(shortTitle string, docs []MarkdownDoc) (document MarkdownDoc, err error) {
	for _, doc := range docs {
		if doc.FrontMatter.Title.Short == shortTitle {
			document = doc
		}
	}
	if document.FrontMatter.Title.Short == "" {
		err = fmt.Errorf("Could not find document with short title: %s", shortTitle)
	}
	return document, err
}

func (d *Docs) FormatFrontMatter(format string, docs []MarkdownDoc) {
	var frontMatterList []DocsFrontMatter
	for _, doc := range docs {
		frontMatterList = append(frontMatterList, doc.FrontMatter)
	}
	switch format {
	case "json":
		fm, _ := json.Marshal(frontMatterList)
		fmt.Print(string(fm))
	default:
		var table strings.Builder
		table.WriteString("| Name | Description | Category | Tags |\n")
		table.WriteString("| ---- | ----------- | -------- | ---- |\n")
		for _, doc := range frontMatterList {
			entry := fmt.Sprintf("| %s | %s | %s | %s |\n", doc.Title.Short, doc.Description, doc.Category, strings.Join(doc.Tags, ", "))
			table.WriteString(entry)
		}
		out, _ := d.Render(table.String())
		fmt.Print(out)
	}
}

func (d *Docs) RenderDocument(doc MarkdownDoc) (string, error) {
	// Add the title since it's captured in frontmatter and not raw markdown
	var bodyWithTitle strings.Builder
	bodyWithTitle.WriteString(fmt.Sprintf("# %s\n", doc.FrontMatter.Title.Long))
	bodyWithTitle.WriteString(doc.Body)
	return d.Render(bodyWithTitle.String())
}
