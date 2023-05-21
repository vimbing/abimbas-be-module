package adidas_backend

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	httpclient "umbrella/internal/http_client"
	definederrors "umbrella/internal/utils/defined_errors"

	http "github.com/vimbing/fhttp"
)

func getEmail(email string) string {
	emailSplit := strings.Split(email, "@")

	if len(emailSplit) < 2 {
		return ""
	}

	prefix := url.QueryEscape(emailSplit[0])

	return fmt.Sprintf("%s@%s", prefix, emailSplit[1])
}

func Login(config *Config) error {
	state := config.TaskStates.Login

	config.DefaultConfig.Log.SetState(state.Name)

	for i := 0; i < state.Retry; i++ {
		config.DefaultConfig.Log.Yellow("Logging in...")
		payload := strings.NewReader(fmt.Sprintf(`grant_type=password&username=%s&password=%s`, getEmail(config.DefaultConfig.TaskData.Email), config.DefaultConfig.TaskData.Password))

		req, err := http.NewRequest("POST", "https://api.3stripes.net/gw-api/v2/token", payload)

		if err != nil {
			config.DefaultConfig.Log.RedDelay(definederrors.MESSAGE_JSON_MARSHALING_ERROR)
			continue
		}

		req.Header = http.Header{
			"host":              {"api.3stripes.net"},
			"accept":            {"application/hal+json"},
			"x-device-info":     {"app/com.adidas.app; os/Android; os-version/29; app-version/5.17.0; buildnumber/51700126; type/willow/Redmi Note 8T/2.75/1080x2130; fingerprint/e4af7e2a081830f9"},
			"x-market":          {"PL"},
			"x-acf-sensor-data": {config.Resources.SensorData},
			"x-api-key":         {"m79qyapn2kbucuv96ednvh22"},
			"accept-language":   {"pl-PL"},
			"user-agent":        {"app/com.adidas.app; os/Android; os-version/29; app-version/5.17.0; buildnumber/51700126; type/willow/Redmi Note 8T/2.75/1080x2130; fingerprint/e4af7e2a081830f9"},
			"x-app-info":        {"platform/Android version/5.17.0"},
			"content-type":      {"application/x-www-form-urlencoded;charset=UTF-8"},
			http.HeaderOrderKey: {
				"host",
				"accept",
				"x-device-info",
				"x-market",
				"x-signature",
				"x-api-key",
				"accept-language",
				"user-agent",
				"x-app-info",
				"x-forter-mobile-uid",
				"content-type",
			},
			http.PHeaderOrderKey: {
				":method",
				":authority",
				":scheme",
				":path",
			},
		}

		res, err := config.DefaultConfig.Network.Client.Do(req)

		if err != nil {
			config.DefaultConfig.Log.RedDelay(definederrors.MESSAGE_REQUEST_SENDING_ERROR)
			continue
		}

		defer res.Body.Close()

		if res.StatusCode != 200 {
			if res.StatusCode == 403 {
				config.DefaultConfig.Log.RedDelay("Error, akamai needs to be solved once again...")
				return AkamaiNeedsSolveError{}
			}

			if res.StatusCode == 401 {
				config.DefaultConfig.Log.RedDelay("Error, account desn't exist...")
				continue
			}

			config.DefaultConfig.Log.StatusCodeErrorDelay(res.Status)
			return WaitAfterError{}
		}

		body, err := httpclient.GetBodyString(res, &config.DefaultConfig.Log)

		if err != nil {
			continue
		}

		var loginResponse LoginResponse

		json.Unmarshal([]byte(body), &loginResponse)

		config.Resources.SessionTokens.AccessToken = loginResponse.AccessToken
		config.Resources.SessionTokens.RefreshToken = loginResponse.RefreshToken

		session, err := GetSessionToSave(config)

		if err != nil {
			return definederrors.ERROR_PLACEHOLDER
		}

		session.Save()

		config.DefaultConfig.Log.Green("Logged in successfully!")

		return nil
	}

	return config.DefaultConfig.Log.LogReturnErrorCustomText(definederrors.ERROR_STOP_TASK, definederrors.MESSAGE_TOO_MANY_RETRYS)
}
