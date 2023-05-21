package adidas_backend

import (
	"strings"
	"umbrella/internal/antibots/bmp/ignasew"

	definederrors "umbrella/internal/utils/defined_errors"
)

func SolveAkamai(config *Config) error {
	state := config.TaskStates.Akamai

	config.DefaultConfig.Log.SetState(state.Name)

	config.DefaultConfig.Log.Yellow("Solving akamai...")

	solverConfig := ignasew.SolverConfig{Website: "adidas"}

	sensor, err := solverConfig.Solve(&config.DefaultConfig)

	if err != nil {
		return err
	}

	config.DefaultConfig.Log.Yellow("Akamai successfully solved...")

	config.Resources.SensorData = strings.Trim(sensor, " ")

	session, err := GetSessionToSave(config)

	if err != nil {
		return definederrors.ERROR_PLACEHOLDER
	}

	session.Save()

	return nil
}
