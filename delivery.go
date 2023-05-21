package adidas_backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	definederrors "umbrella/internal/utils/defined_errors"

	http "github.com/vimbing/fhttp"
)

func Delivery(config *Config) error {
	state := config.TaskStates.Delivery

	config.DefaultConfig.Log.SetState(state.Name)

	for i := 0; i < state.Retry; i++ {
		config.DefaultConfig.Log.Yellow("Filling delivery data...")

		deliveryPayload := DeliveryPayload{
			ShippingMethodID: "Standard-EFC-3",
		}

		deliveryPayloadJson, err := json.Marshal(deliveryPayload)

		if err != nil {
			config.DefaultConfig.Log.RedDelay(definederrors.MESSAGE_JSON_MARSHALING_ERROR)
			continue
		}

		req, err := http.NewRequest("PUT", fmt.Sprintf("https://api.3stripes.net/gw-api/v2/checkouts/%s/lines/me/shipping_method", config.Resources.CheckoutID), bytes.NewBuffer(deliveryPayloadJson))

		if err != nil {
			config.DefaultConfig.Log.RedDelay(definederrors.MESSAGE_JSON_MARSHALING_ERROR)
			continue
		}

		req.Header = http.Header{
			"host":            {"api.3stripes.net"},
			"x-app-info":      {"platform/iOSversion/5.17"},
			"user-agent":      {"adidas/2022.11.25.13.12CFNetwork/1390Darwin/22.0.0"},
			"x-device-info":   {"app/adidas;os/iOS;os-version/16.0;app-version/5.17;buildnumber/2022.11.25.13.12;type/iPhone12,1;fingerprint/8AFED4D6-A075-4108-9965-BD2821D7F736"},
			"x-market":        {"PL"},
			"authorization":   {config.Resources.SessionTokens.AccessToken},
			"accept-language": {"pl-PL"},
			"accept":          {"application/hal+json"},
			"content-type":    {"application/json;charset=UTF-8"},
			"x-api-key":       {"m79qyapn2kbucuv96ednvh22"},
			http.HeaderOrderKey: {
				"host",
				"x-product-cache",
				"x-app-info",
				"user-agent",
				"x-pdata",
				"x-feed-cache",
				"x-device-info",
				"x-forter-mobile-uid",
				"x-pdata-cache",
				"x-signature",
				"x-market",
				"authorization",
				"accept-language",
				"accept",
				"content-type",
				"x-api-key",
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
			return WaitAfterError{}
		}

		config.DefaultConfig.Log.Green("Delivery data filled successfully!")

		return nil
	}

	return config.DefaultConfig.Log.LogReturnErrorCustomText(definederrors.ERROR_TOO_MANY_RETRYS, definederrors.MESSAGE_TOO_MANY_RETRYS)
}
