package adidas_backend

import (
	"bytes"
	"encoding/json"
	profilesreader "umbrella/internal/file_readers/profiles_reader"
	httpclient "umbrella/internal/http_client"
	definederrors "umbrella/internal/utils/defined_errors"

	http "github.com/vimbing/fhttp"
)

func Register(config *Config) error {
	state := config.TaskStates.Register

	config.DefaultConfig.Log.SetState(state.Name)

	for i := 0; i < state.Retry; i++ {
		config.DefaultConfig.Log.Yellow("Creating account...")

		registerPayload := RegisterPayload{
			Email:             profilesreader.HandleEmailJig(config.DefaultConfig.TaskData.Email),
			Password:          config.DefaultConfig.TaskData.Password,
			MembershipConsent: true,
			DormantPeriod:     "1y",
		}

		registerPayloadJson, err := json.Marshal(registerPayload)

		if err != nil {
			config.DefaultConfig.Log.Yellow(definederrors.MESSAGE_JSON_MARSHALING_ERROR)
			continue
		}

		req, err := http.NewRequest("POST", "https://api.3stripes.net/gw-api/v2/user", bytes.NewBuffer(registerPayloadJson))

		if err != nil {
			config.DefaultConfig.Log.RedDelay(definederrors.MESSAGE_JSON_MARSHALING_ERROR)
			continue
		}

		req.Header = http.Header{
			"host":                {"api.3stripes.net"},
			"accept":              {"application/hal+json"},
			"x-device-info":       {"app/com.adidas.app; os/Android; os-version/29; app-version/5.17.0; buildnumber/51700126; type/willow/Redmi Note 8T/2.75/1080x2130; fingerprint/e4af7e2a081830f9"},
			"x-market":            {"PL"},
			"x-acf-sensor-data":   {config.Resources.SensorData},
			"x-api-key":           {"m79qyapn2kbucuv96ednvh22"},
			"accept-language":     {"pl-PL"},
			"user-agent":          {"app/com.adidas.app; os/Android; os-version/29; app-version/5.17.0; buildnumber/51700126; type/willow/Redmi Note 8T/2.75/1080x2130; fingerprint/e4af7e2a081830f9"},
			"x-app-info":          {"platform/Android version/5.17.0"},
			"x-forter-mobile-uid": {"8AFED4D6-A075-4108-9965-BD2821D7F736"},
			"content-type":        {"application/json;charset=UTF-8"},
			http.HeaderOrderKey: {
				"host",
				"accept",
				"x-device-info",
				"x-market",
				"x-signature",
				"x-acf-sensor-data",
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
				config.DefaultConfig.Log.Yellow("Akamai needs to be solved...")
				return AkamaiNeedsSolveError{}
			}

			config.DefaultConfig.Log.StatusCodeErrorDelay(res.Status)
			return definederrors.ERROR_STOP_TASK
		}

		body, err := httpclient.GetBodyString(res, &config.DefaultConfig.Log)

		if err != nil {
			continue
		}

		var registerResponse RegisterResponse

		json.Unmarshal([]byte(body), &registerResponse)

		config.Resources.SessionTokens.AccessToken = registerResponse.AccessToken
		config.Resources.SessionTokens.RefreshToken = registerResponse.RefreshToken

		session, err := GetSessionToSave(config)

		if err != nil {
			return definederrors.ERROR_PLACEHOLDER
		}

		session.Save()

		config.DefaultConfig.Log.Green("Account created successfully!")

		return nil
	}

	return config.DefaultConfig.Log.LogReturnErrorCustomText(definederrors.ERROR_TOO_MANY_RETRYS, definederrors.MESSAGE_TOO_MANY_RETRYS)
}
