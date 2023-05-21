package adidas_backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	definederrors "umbrella/internal/utils/defined_errors"
	discountmenager "umbrella/internal/utils/discount_menager"

	http "github.com/vimbing/fhttp"
)

func Discount(config *Config) error {
	if len(config.DefaultConfig.TaskData.Discount) < 3 {
		return nil
	}

	state := config.TaskStates.Discount

	config.DefaultConfig.Log.SetState(state.Name)

	for i := 0; i < state.Retry; i++ {
		config.DefaultConfig.Log.Yellow("Adding discount code...")

		code, err := discountmenager.GetCode(config.DefaultConfig.TaskData.Discount, config.DefaultConfig.TaskId)

		if err != nil {
			config.DefaultConfig.Log.RedDelay("No available discount codes in file!")
			return definederrors.ERROR_STOP_TASK
		}

		config.UsedDiscount = code
		config.IfDiscounted = true

		discountPayload := DiscountPayload{
			VoucherCode: code,
		}

		discountPayloadJson, err := json.Marshal(discountPayload)

		if err != nil {
			config.DefaultConfig.Log.RedDelay(definederrors.MESSAGE_JSON_MARSHALING_ERROR)
			continue
		}

		req, err := http.NewRequest("PUT", fmt.Sprintf("https://api.3stripes.net/gw-api/v2/checkouts/%s/voucher", config.Resources.CheckoutID), bytes.NewBuffer(discountPayloadJson))

		if err != nil {
			config.DefaultConfig.Log.RedDelay(definederrors.MESSAGE_JSON_MARSHALING_ERROR)
			continue
		}

		req.Header = http.Header{
			"host":            {"api.3stripes.net"},
			"x-app-info":      {"platform/Android version/5.17.0"},
			"user-agent":      {"app/com.adidas.app; os/Android; os-version/29; app-version/5.17.0; buildnumber/51700126; type/willow/Redmi Note 8T/2.75/1080x2130; fingerprint/e4af7e2a081830f9"},
			"x-device-info":   {"app/com.adidas.app; os/Android; os-version/29; app-version/5.17.0; buildnumber/51700126; type/willow/Redmi Note 8T/2.75/1080x2130; fingerprint/e4af7e2a081830f9"},
			"authorization":   {config.Resources.SessionTokens.AccessToken},
			"x-signature":     {"33A97E353F88790F00AAF989951CBA4F0F8225463661F9B5DFEFC8799E692982"},
			"x-market":        {"PL"},
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
				"x-acf-sensor-data",
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

		config.DefaultConfig.Log.Green("Discount code added successfully!")

		return nil
	}

	return config.DefaultConfig.Log.LogReturnErrorCustomText(definederrors.ERROR_TOO_MANY_RETRYS, definederrors.MESSAGE_TOO_MANY_RETRYS)
}
