package adidas_backend

import (
	definederrors "umbrella/internal/utils/defined_errors"

	http "github.com/vimbing/fhttp"
)

func CheckSession(config *Config) error {
	state := config.TaskStates.SessionCheck

	config.DefaultConfig.Log.SetState(state.Name)

	for i := 0; i < state.Retry; i++ {
		config.DefaultConfig.Log.Yellow("Checking session...")

		req, err := http.NewRequest("GET", "https://api.3stripes.net/gw-api/v2/user", nil)

		if err != nil {
			config.DefaultConfig.Log.RedDelay(definederrors.MESSAGE_JSON_MARSHALING_ERROR)
			continue
		}

		req.Header = http.Header{
			"host":                {"api.3stripes.net"},
			"accept":              {"application/hal+json"},
			"x-device-info":       {"app/adidas;os/iOS;os-version/16.0;app-version/5.17;buildnumber/2022.11.25.13.12;type/iPhone12,1;fingerprint/8AFED4D6-A075-4108-9965-BD2821D7F736"},
			"x-market":            {"PL"},
			"authorization":       {config.Resources.SessionTokens.AccessToken},
			"x-signature":         {"0A68DE58857F0EF229BA5C0F81387E6EB02F06AF7671FE0D674E3622192C6B7D"},
			"accept-language":     {"pl-PL"},
			"x-api-key":           {"m79qyapn2kbucuv96ednvh22"},
			"user-agent":          {"adidas/2022.11.25.13.12CFNetwork/1390Darwin/22.0.0"},
			"x-app-info":          {"platform/iOSversion/5.17"},
			"x-forter-mobile-uid": {"8AFED4D6-A075-4108-9965-BD2821D7F736"},
			http.HeaderOrderKey: {
				"host",
				"accept",
				"x-device-info",
				"x-market",
				"authorization",
				"x-signature",
				"x-pdata-cache",
				"accept-language",
				"x-api-key",
				"x-feed-cache",
				"x-product-cache",
				"user-agent",
				"x-app-info",
				"x-forter-mobile-uid",
				"x-pdata",
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
			config.DefaultConfig.Log.Yellow("Session is not valid, proceeding to login...")
			return definederrors.ERROR_PLACEHOLDER
		}

		config.DefaultConfig.Log.Green("Saved session is valid, skipping login!")

		return nil
	}

	return config.DefaultConfig.Log.LogReturnErrorCustomText(definederrors.ERROR_TOO_MANY_RETRYS, definederrors.MESSAGE_TOO_MANY_RETRYS)
}
