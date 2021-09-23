package run

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	containerId string
)

func CreateRunCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "run",
		Short: "Run a validator container against a given resource dir",
		Long:  `Run a validator container against a given resource dir`,
		Run:   runValidator,
	}

	tmp.Flags().StringVar(&containerId, "validator-container-id", "", "The PCV validator container to run")

	return tmp
}

func runValidator(cmd *cobra.Command, args []string) {

	resourceDirPath, err := filepath.Abs(args[0])
	if err != nil {
		panic("Could not generate path to resource")
	}

	/* #nosec G204 */
	runCommand := exec.Command("docker", "run", "-v", fmt.Sprintf("%v:/module", resourceDirPath), "-w", "/module", containerId)

	var out bytes.Buffer
	var stderr bytes.Buffer

	runCommand.Stdout = &out
	runCommand.Stderr = &stderr

	err = runCommand.Run()

	log.Error().Msgf(stderr.String())
	log.Info().Msgf(out.String())

	if err != nil {
		log.Fatal().Msgf(err.Error())
	}
}
