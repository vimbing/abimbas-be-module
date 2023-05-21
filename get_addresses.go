package adidas_backend

import (
	"encoding/json"
	httpclient "umbrella/internal/http_client"
	definederrors "umbrella/internal/utils/defined_errors"

	http "github.com/vimbing/fhttp"
)

func GetAddresses(config *Config) error {
	state := config.TaskStates.AddressGet

	config.DefaultConfig.Log.SetState(state.Name)

	for i := 0; i < state.Retry; i++ {
		config.DefaultConfig.Log.Yellow("Getting addresses from account...")

		req, err := http.NewRequest("GET", "https://api.3stripes.net/gw-api/v2/user/addresses", nil)

		if err != nil {
			config.DefaultConfig.Log.RedDelay(definederrors.MESSAGE_JSON_MARSHALING_ERROR)
			continue
		}

		req.Header = http.Header{
			"host":              {"api.3stripes.net"},
			"accept":            {"application/hal+json"},
			"x-device-info":     {"app/com.adidas.app; os/Android; os-version/29; app-version/5.17.0; buildnumber/51700126; type/willow/Redmi Note 8T/2.75/1080x2130; fingerprint/e4af7e2a081830f9"},
			"x-market":          {"PL"},
			"authorization":     {config.Resources.SessionTokens.AccessToken},
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
				"authorization",
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
			config.DefaultConfig.Log.StatusCodeErrorDelay(res.Status)
			return definederrors.ERROR_STOP_TASK
		}

		body, err := httpclient.GetBodyString(res, &config.DefaultConfig.Log)

		if err != nil {
			continue
		}

		var getAddressesResponse GetAddressesResponse

		json.Unmarshal([]byte(body), &getAddressesResponse)

		if len(getAddressesResponse.Addresses) < 1 {
			config.DefaultConfig.Log.Green("No available addresses for this account!")
			return definederrors.ERROR_STOP_TASK
		}

		config.Resources.AddressId = getAddressesResponse.Addresses[0].ID

		config.DefaultConfig.Log.Green("Address successfully scraped!")

		return nil
	}

	return config.DefaultConfig.Log.LogReturnErrorCustomText(definederrors.ERROR_TOO_MANY_RETRYS, definederrors.MESSAGE_TOO_MANY_RETRYS)
}
