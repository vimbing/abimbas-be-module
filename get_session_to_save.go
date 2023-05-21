package adidas_backend

import (
	"encoding/json"
	"strings"
	customjar "umbrella/internal/http_client/custom_jar"
	sessionsaver "umbrella/internal/session_saver"
	definederrors "umbrella/internal/utils/defined_errors"
)

func GetSessionToSave(config *Config) (sessionsaver.Session, error) {
	loginResponse := LoginResponse{
		AccessToken:  config.Resources.SessionTokens.AccessToken,
		RefreshToken: config.Resources.SessionTokens.RefreshToken,
	}

	loginResponseJson, err := json.Marshal(loginResponse)

	if err != nil {
		config.DefaultConfig.Log.RedDelay(definederrors.MESSAGE_JSON_MARSHALING_ERROR)
		return sessionsaver.Session{}, err
	}

	session := sessionsaver.Session{
		Website: strings.ToLower(config.DefaultConfig.TaskData.Website),
		Email:   strings.ToLower(config.DefaultConfig.TaskData.Email),
		Cookies: []customjar.Cookie{
			{
				Value: string(loginResponseJson),
				Name:  "body",
			},
			GetSensorAsCookie(config),
		},
		UserAgent: config.DefaultConfig.Network.UserAgent,
		Proxy:     config.DefaultConfig.Network.Proxy,
	}

	return session, nil
}
