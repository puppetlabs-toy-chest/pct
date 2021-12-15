package explain

import (
	"fmt"

	"github.com/puppetlabs/pdkgo/docs/md"
	"github.com/puppetlabs/pdkgo/pkg/docs"
	"github.com/spf13/cobra"
)

var (
	docsApi    *docs.Docs
	listTopics bool
	format     string
	tag        string
	category   string
	topic      string
	// Possibly implement later to enable context aware filtering for tags/categories
	// depending on which is filtered first?
	// filteredDocs []docs.MarkdownDoc
)

func CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "explain",
		Short:             "Present documentation about topics",
		Long:              "Present documentation about various topics, including...",
		Args:              validateArgCount,
		ValidArgsFunction: flagCompletion,
		PreRun:            preExecute,
		RunE:              execute,
	}

	dfs := md.GetDocsFS()
	docsApi = &docs.Docs{
		DocsFileSystem: &dfs,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().BoolVarP(&listTopics, "list", "l", false, "list available topics")

	cmd.Flags().StringVarP(&format, "format", "f", "human", "display output in human or json format")
	err := cmd.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"table", "json"}, cobra.ShellCompDirectiveNoFileComp
	})
	cobra.CheckErr(err)

	cmd.Flags().StringVarP(&tag, "tag", "t", "", "filter available topics by tag")
	err = cmd.RegisterFlagCompletionFunc("tag", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if docsApi.ParsedDocsCache == nil {
			preExecute(cmd, args)
		}
		return docsApi.ListTags(docsApi.ParsedDocsCache), cobra.ShellCompDirectiveNoFileComp
	})
	cobra.CheckErr(err)

	cmd.Flags().StringVarP(&category, "category", "c", "", "filter available topics by category")
	err = cmd.RegisterFlagCompletionFunc("category", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if docsApi.ParsedDocsCache == nil {
			preExecute(cmd, args)
		}
		return docsApi.ListCategories(docsApi.ParsedDocsCache), cobra.ShellCompDirectiveNoFileComp
	})
	cobra.CheckErr(err)

	return cmd
}

func preExecute(cmd *cobra.Command, args []string) {
	docsApi.FindAndParse("content")
}

func validateArgCount(cmd *cobra.Command, args []string) error {
	// show available topics if user runs `pct explain`
	if len(args) == 0 && !listTopics {
		listTopics = true
	}

	if len(args) == 1 {
		if category != "" || tag != "" {
			return fmt.Errorf("Specify a topic *or* search by tag/category")
		}
		topic = args[0]
	} else if len(args) > 1 {
		return fmt.Errorf("Specify only one topic to explain")
	}

	return nil
}

func flagCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if docsApi.ParsedDocsCache == nil {
		preExecute(cmd, args)
	}
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return docsApi.CompleteTitle(docsApi.ParsedDocsCache, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
}

func execute(cmd *cobra.Command, args []string) error {
	docs := docsApi.ParsedDocsCache
	if listTopics {
		if format == "" {
			format = "table"
		}
		if category != "" {
			docs = docsApi.FilterByCategory(category, docs)
		}
		if tag != "" {
			docs = docsApi.FilterByTag(tag, docs)
		}
		// If there's only one match, should we render the matching doc?
		docsApi.FormatFrontMatter(format, docs)
	} else if topic != "" {
		doc, err := docsApi.SelectDocument(topic, docsApi.ParsedDocsCache)
		if err != nil {
			return err
		}
		output, err := docsApi.RenderDocument(doc)
		if err != nil {
			return err
		}
		fmt.Print(output)
		// If --online, open in browser and do not display
		// Should we have a --scroll mode?
	}
	return nil
}
