package adidas_backend

import (
	"bytes"
	"encoding/json"
	httpclient "umbrella/internal/http_client"
	clititle "umbrella/internal/utils/cli_title"
	definederrors "umbrella/internal/utils/defined_errors"

	http "github.com/vimbing/fhttp"
)

func CheckoutId(config *Config, variant Variant) error {
	state := config.TaskStates.CheckoutID

	config.DefaultConfig.Log.SetState(state.Name)

	for i := 0; i < state.Retry; i++ {
		config.DefaultConfig.Log.Yellow("Getting checkout id...")

		checkoutIdPayload := CheckoutIdPayload{
			Items: []Items{
				{
					ProductID:          variant.Pid,
					VariationProductID: variant.SizePid,
					Quantity:           config.DefaultConfig.TaskData.Quantity.GetQuantity(),
				},
			},
		}

		checkoutIdPayloadJson, err := json.Marshal(checkoutIdPayload)

		if err != nil {
			config.DefaultConfig.Log.RedDelay(definederrors.MESSAGE_JSON_MARSHALING_ERROR)
			continue
		}

		req, err := http.NewRequest("POST", "https://api.3stripes.net/gw-api/v2/checkouts", bytes.NewBuffer(checkoutIdPayloadJson))

		if err != nil {
			config.DefaultConfig.Log.RedDelay(definederrors.MESSAGE_JSON_MARSHALING_ERROR)
			continue
		}

		req.Header = http.Header{
			"host":              {"api.3stripes.net"},
			"x-device-info":     {"app/com.adidas.app; os/Android; os-version/29; app-version/5.17.0; buildnumber/51700126; type/willow/Redmi Note 8T/2.75/1080x2130; fingerprint/e4af7e2a081830f9"},
			"user-agent":        {"app/com.adidas.app; os/Android; os-version/29; app-version/5.17.0; buildnumber/51700126; type/willow/Redmi Note 8T/2.75/1080x2130; fingerprint/e4af7e2a081830f9"},
			"x-app-info":        {"platform/Android version/5.17.0"},
			"x-market":          {"PL"},
			"accept-language":   {"pl-PL"},
			"authorization":     {config.Resources.SessionTokens.AccessToken},
			"x-acf-sensor-data": {config.Resources.SensorData},
			"accept":            {"application/hal+json"},
			"content-type":      {"application/json;charset=UTF-8"},
			"x-api-key":         {"m79qyapn2kbucuv96ednvh22"},
			http.HeaderOrderKey: {
				"host",
				"x-product-cache",
				"x-app-info",
				"user-agent",
				"x-pdata",
				"authorization",
				"x-feed-cache",
				"x-device-info",
				"x-forter-mobile-uid",
				"x-pdata-cache",
				"x-signature",
				"x-market",
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
			if res.StatusCode == 403 {
				config.DefaultConfig.Log.RedDelay("Error, akamai needs to be solved once again...")
				return AkamaiNeedsSolveError{}
			}

			if res.StatusCode == 400 {
				config.DefaultConfig.Log.StatusCodeErrorDelay(res.Status)
				return WaitAfterRestock{}
			}

			config.DefaultConfig.Log.StatusCodeErrorDelay(res.Status)
			return WaitAfterError{}
		}

		body, err := httpclient.GetBodyString(res, &config.DefaultConfig.Log)

		if err != nil {
			continue
		}

		var checkoutIdResponse CheckoutIdResponse

		json.Unmarshal([]byte(body), &checkoutIdResponse)

		config.Resources.CheckoutID = checkoutIdResponse.ID

		if checkoutIdResponse.Selected.BillingAddress.ID != config.DefaultConfig.Profile.Payment {
			checkoutIdResponse.Selected.PaymentMethod.ID = ""
		}

		config.Cosmetics = GetCosmetics(&checkoutIdResponse)

		config.Resources.RequiredRequests = RequiredRequests{
			Shipping: !(len(checkoutIdResponse.Selected.BillingAddress.ID) > 0),
			Payment:  !(len(checkoutIdResponse.Selected.PaymentMethod.ID) > 0),
		}

		config.DefaultConfig.Log.Green("Checkout id successfully scraped!")

		clititle.AddCarted()

		return nil
	}

	return config.DefaultConfig.Log.LogReturnErrorCustomText(definederrors.ERROR_TOO_MANY_RETRYS, definederrors.MESSAGE_TOO_MANY_RETRYS)
}
