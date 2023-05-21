package adidas_backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	adyen "umbrella/internal/antibots/adyen"
	httpclient "umbrella/internal/http_client"
	definederrors "umbrella/internal/utils/defined_errors"

	http "github.com/vimbing/fhttp"
)

func OrderCreditCard(config *Config) error {
	state := config.TaskStates.Order

	config.DefaultConfig.Log.SetState(state.Name)

	for i := 0; i < state.Retry; i++ {
		config.DefaultConfig.Log.Yellow("Creating order...")

		clientSideEncrypter := adyen.ClientSideEncrypter{AdyenPublicKey: ADYEN_PUBLIC}
		adyenNonce, err := clientSideEncrypter.GenerateAdyenNonce(config.DefaultConfig.Profile.FirstName, config.DefaultConfig.Profile.CardNumber, config.DefaultConfig.Profile.Cvv, config.DefaultConfig.Profile.Month, config.DefaultConfig.Profile.Year)

		if err != nil {
			return err
		}

		orderPayloadJson, err := json.Marshal(PaymentCCPayload{
			CheckoutID: config.Resources.CheckoutID,
			NewCard: NewCard{
				Type:       strings.ToUpper(config.DefaultConfig.Profile.CardType),
				Encryption: adyenNonce,
			},
		})

		fmt.Println(string(orderPayloadJson))

		if err != nil {
			config.DefaultConfig.Log.RedDelay(definederrors.MESSAGE_JSON_MARSHALING_ERROR)
			continue
		}

		req, err := http.NewRequest("POST", "https://api.3stripes.net/gw-api/v2/orders", bytes.NewBuffer(orderPayloadJson))

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
			"authorization":     {config.Resources.SessionTokens.AccessToken},
			"accept-language":   {"pl-PL"},
			"accept":            {"application/hal+json"},
			"content-type":      {"application/json;charset=UTF-8"},
			"x-api-key":         {"m79qyapn2kbucuv96ednvh22"},
			"x-acf-sensor-data": {config.Resources.SensorData},
			http.HeaderOrderKey: {
				"host",
				"x-product-cache",
				"x-app-info",
				"cookie",
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

		body, err := httpclient.GetBodyString(res, &config.DefaultConfig.Log)

		if err != nil {
			continue
		}

		fmt.Println(body)
		os.Exit(1)

		defer res.Body.Close()

		if res.StatusCode != 200 {
			if res.StatusCode == 409 && strings.Contains(body, "is not supported") {
				config.DefaultConfig.Log.Yellow("COD is not supported, switching to paypal")
				config.DefaultConfig.Profile.Payment = PAYMENT_PAYPAL
				continue
			}

			config.DefaultConfig.Log.StatusCodeErrorDelay(res.Status)
			return WaitAfterError{}
		}

		config.DefaultConfig.Log.Green("Payment data filled successfully!")

		return nil
	}

	return config.DefaultConfig.Log.LogReturnErrorCustomText(definederrors.ERROR_TOO_MANY_RETRYS, definederrors.MESSAGE_TOO_MANY_RETRYS)
}
