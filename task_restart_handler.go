package adidas_backend

import (
	"errors"
	"time"
)

func HandleTaskRestart(config *Config, err error) error {
	if errors.Is(AkamaiNeedsSolveError{}, err) {
		err = SolveAkamai(config)

		if err != nil {
			return err
		}
	}

	if errors.Is(WaitAfterRestock{}, err) {
		config.DefaultConfig.Log.Yellow("Fake stock presented, waiting for refresh...")
		time.Sleep(time.Second * 45)
	}

	if errors.Is(WaitAfterError{}, err) {
		config.DefaultConfig.Log.Yellow("Error found, waiting one minute...")
		time.Sleep(time.Minute * 1)
	}

	if errors.Is(RefreshSessionError{}, err) {
		err = RefreshToken(config)

		if err != nil {
			return err
		}
	}

	return nil
}
