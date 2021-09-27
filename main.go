package main

import (
	"context"
	"net/http"

	"github.com/puppetlabs/pdkgo/cmd/build"
	"github.com/puppetlabs/pdkgo/cmd/completion"
	"github.com/puppetlabs/pdkgo/cmd/install"
	"github.com/puppetlabs/pdkgo/cmd/new"
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

	// Context must always be established:
	ctx := context.Background()
	telemetry.Start(ctx, honeycomb_api_key, honeycomb_dataset)

	var rootCmd = root.CreateRootCommand()

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
		},
	}
	rootCmd.AddCommand(installCmd.CreateCommand())

	// new
	rootCmd.AddCommand(new.CreateCommand())

	cobra.OnInitialize(root.InitLogger, root.InitConfig)
	cobra.CheckErr(rootCmd.Execute())
}
