package adidas_backend

import (
	"encoding/json"
	"errors"
	"fmt"
	sessionsaver "umbrella/internal/session_saver"
)

func LoadSession(session sessionsaver.Session, config *Config) error {
	if len(session.Cookies) < 1 {
		return errors.New("placeholder")
	}

	var loginResponse LoginResponse

	for _, c := range session.Cookies {
		if c.Name == "body" {
			json.Unmarshal([]byte(c.Value), &loginResponse)
		}
		if c.Name == "sensor_data" {
			config.Resources.SensorData = c.Value
		}
	}

	config.Resources.SessionTokens.AccessToken = fmt.Sprintf("Bearer %s", loginResponse.AccessToken)
	config.Resources.SessionTokens.RefreshToken = loginResponse.RefreshToken

	return nil
}

func SavedSessionHandler(config *Config) error {
	config.DefaultConfig.Log.SetState(config.TaskStates.SessionCheck.Name)

	session, err := sessionsaver.GetSession(config.DefaultConfig.TaskData.Email, config.DefaultConfig.TaskData.Website)

	if err == nil {
		config.DefaultConfig.Log.Yellow("Found saved session!")

		err = LoadSession(session, config)

		if err != nil {
			return err
		}

		err = RefreshToken(config)

		if err == nil {
			return nil
		}
	}

	err = SolveAkamai(config)

	if err != nil {
		return err
	}

	err = Login(config)

	if err != nil {
		return err
	}

	return nil
}
