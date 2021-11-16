package main

import (
	"context"
	"net/http"

	"github.com/puppetlabs/pdkgo/internal/pkg/exec_runner"

	"github.com/puppetlabs/pdkgo/cmd/build"
	"github.com/puppetlabs/pdkgo/cmd/completion"
	"github.com/puppetlabs/pdkgo/cmd/install"
	"github.com/puppetlabs/pdkgo/cmd/new"
	"github.com/puppetlabs/pdkgo/cmd/news"
	"github.com/puppetlabs/pdkgo/cmd/root"
	appver "github.com/puppetlabs/pdkgo/cmd/version"
	"github.com/puppetlabs/pdkgo/internal/pkg/gzip"
	"github.com/puppetlabs/pdkgo/internal/pkg/pct"
	"github.com/puppetlabs/pdkgo/internal/pkg/tar"
	"github.com/puppetlabs/pdkgo/pkg/telemetry"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	version           = "dev"
	commit            = "none"
	date              = "unknown"
	honeycomb_api_key = "not_set"
	honeycomb_dataset = "not_set"
)

func main() {
	// Telemetry must be initialized before anything else;
	// If the telemetry build tag was not passed, this is all null ops
	ctx, traceProvider, parentSpan := telemetry.Start(context.Background(), honeycomb_api_key, honeycomb_dataset, "pct")

	var rootCmd = root.CreateRootCommand()

	// Get the command called and its arguments;
	// The arguments are only necessary if we want to
	// hand them off as an attribute to the parent span:
	// do we? Otherwise we just need the calledCommand
	calledCommand, calledCommandArguments := root.GetCalledCommand(rootCmd)
	telemetry.AddStringSpanAttribute(parentSpan, "arguments", calledCommandArguments)

	var verCmd = appver.CreateVersionCommand(version, date, commit)
	v := appver.Format(version, date, commit)
	rootCmd.Version = v
	rootCmd.SetVersionTemplate(v)
	rootCmd.AddCommand(verCmd)

	rootCmd.AddCommand(completion.CreateCompletionCommand())

	// afero setup
	fs := afero.NewOsFs()
	afs := afero.Afero{Fs: fs}
	iofs := afero.IOFS{Fs: fs}

	// build
	rootCmd.AddCommand(build.CreateCommand())

	// install
	installCmd := install.InstallCommand{
		PctInstaller: &pct.PctInstaller{
			Tar:        &tar.Tar{AFS: &afs},
			Gunzip:     &gzip.Gunzip{AFS: &afs},
			AFS:        &afs,
			IOFS:       &iofs,
			HTTPClient: &http.Client{},
			Exec:       &exec_runner.Exec{},
		},
		AFS: &afs,
	}
	rootCmd.AddCommand(installCmd.CreateCommand())

	// new
	rootCmd.AddCommand(new.CreateCommand())

	// news
	rootCmd.AddCommand(news.CreateCommand())

	// initialize
	cobra.OnInitialize(root.InitLogger, root.InitConfig)

	// instrument & execute called command
	ctx, childSpan := telemetry.NewSpan(ctx, calledCommand)
	err := rootCmd.ExecuteContext(ctx)
	telemetry.RecordSpanError(childSpan, err)
	telemetry.EndSpan(childSpan)

	// Send all events
	telemetry.ShutDown(ctx, traceProvider, parentSpan)

	// Handle exiting with/out errors.
	cobra.CheckErr(err)
}
