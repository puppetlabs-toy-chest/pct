package puppet

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/puppetlabs/pdkgo/internal/pkg/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)


func StartContainer(puppetVersion int) {
	containerName := fmt.Sprintf("pdk:puppet%d", puppetVersion)
	log.Info().Msgf("Starting container `%s`", containerName)

	runCommand := exec.Command("docker", "run", "--rm", "-i", "-d", "--name", "pdk-toolchain", "-v", "/toolchain", containerName)

	if err := runCommand.Run(); err != nil {
		log.Error().Msgf("Unable to start the toolchain: %s", err.Error())
	}

	storeMount()
}

func StopContainer() {
	log.Info().Msgf("Stopping toolchain container")

	runCommand := exec.Command("docker", "stop", "pdk-toolchain")

	if err := runCommand.Run(); err != nil {
		if exiterr, _ := err.(*exec.ExitError); exiterr.ExitCode() > 1 {
			log.Error().Msgf("Unable to stop the toolchain container: %s", exiterr.Error())
		}
	}

}

func storeMount() {
	runCommand, err := exec.Command("docker", "container", "inspect", "-f '{{ (index .Mounts 0).Source }}'", "pdk-toolchain").Output()

	if err != nil {
		log.Error().Msgf("Failed to record toolchain mount location: %s", err.Error())
		return
	}

	mountPath := strings.Trim(string(runCommand), "\n' ")
	log.Info().Msgf("mount location: %s", mountPath)

	viper.Set(config.ToolchainMount, mountPath)
	if err:= viper.WriteConfig(); err != nil {
		log.Error().Msgf("Failed to save mountpath to config: %s", err.Error())
	}
}


func GetToolchainPath() string {
	return viper.GetString(config.ToolchainMount)
}
