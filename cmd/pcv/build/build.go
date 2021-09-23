package build

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	validatorDirPath string
	dockerfile       string
)

type PuppetValidatorTemplateInfo struct {
	ID     string `yaml:"id"`
	Author string `yaml:"author"`
}

func CreateBuildCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "build",
		Short: "Build a validator container from a Dockerfile",
		Long:  `Build a validator container from a Dockerfile`,
		Run:   build,
	}

	tmp.Flags().StringVar(&validatorDirPath, "validator-template", "", "The path to the PCV template")

	return tmp
}

func getValidatorConfig() (info PuppetValidatorTemplateInfo) {
	pcvConfigFilePath := filepath.Join(validatorDirPath, "pcv-config.yml")
	if _, err := os.Stat(pcvConfigFilePath); err != nil {
		log.Fatal().Msgf("No pcv-config.yml found at: %v", pcvConfigFilePath)
		return
	}
	fileBytes, err := ioutil.ReadFile(pcvConfigFilePath)
	if err != nil {
		log.Fatal().Msgf("Error reading %v", pcvConfigFilePath)
		return
	}

	err = yaml.Unmarshal(fileBytes, &info)
	if err != nil {
		panic("Could not unmarshal pcv-config.yml")
	}

	return info
}

func buildContainer() {
	log.Info().Msgf("Building container from Dockerfile: \n%s", dockerfile)

	info := getValidatorConfig()

	dockerfileFilePath := filepath.Join(validatorDirPath, "content")
	if _, err := os.Stat(dockerfileFilePath); err != nil {
		log.Fatal().Msgf("No Dockerfile found at: %v", dockerfileFilePath)
		return
	}

	/* #nosec G204 */
	runCommand := exec.Command("docker", "build", dockerfileFilePath, "--rm", "--tag", info.ID)

	err := runCommand.Run()

	if err != nil {
		log.Fatal().Msgf(err.Error())
	}
}

func build(cmd *cobra.Command, args []string) {
	buildContainer()
}
