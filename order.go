package adidas_backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	httpclient "umbrella/internal/http_client"
	"umbrella/internal/utils/consts"
	definederrors "umbrella/internal/utils/defined_errors"
	discountmenager "umbrella/internal/utils/discount_menager"
	successhandler "umbrella/internal/utils/success_handler"
	webhookenginev2 "umbrella/internal/webhook_engine_v2"

	http "github.com/vimbing/fhttp"
)

func Order(config *Config) error {
	state := config.TaskStates.Order

	config.DefaultConfig.Log.SetState(state.Name)

	for i := 0; i < state.Retry; i++ {
		config.DefaultConfig.Log.Yellow("Creating order...")

		orderPayload := OrderPayload{
			CheckoutID: config.Resources.CheckoutID,
		}

		orderPayloadJson, err := json.Marshal(orderPayload)

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
			"x-product-cache":   {"1535"},
			"x-device-info":     {"app/com.adidas.app; os/Android; os-version/29; app-version/5.17.0; buildnumber/51700126; type/willow/Redmi Note 8T/2.75/1080x2130; fingerprint/e4af7e2a081830f9"},
			"user-agent":        {"app/com.adidas.app; os/Android; os-version/29; app-version/5.17.0; buildnumber/51700126; type/willow/Redmi Note 8T/2.75/1080x2130; fingerprint/e4af7e2a081830f9"},
			"x-app-info":        {"platform/Android version/5.17.0"},
			"x-pdata":           {"H4sIAAAAAAAAA6tWykxRslKKsHA2dnI1D/H1cg0ycAsJUdJRyk3NTUotUrKqhqhwdPEM8DEyNbIwNjQ3MrCwAKooyQTJG+goFaWWJxalFCtZRRtZ6hia6xjpGJnE1uoolebl5Cdnp6a4pSaWlBalghTE6iglJ+YWJGam54G5tQAOHVSmgQAAAA=="},
			"x-feed-cache":      {"-1634373788"},
			"x-pdata-cache":     {"-1634373788"},
			"x-signature":       {"7FB5017D81718594D9294AAD85AB9BCDC516A501377820A61A3026D90A872606"},
			"x-market":          {"PL"},
			"authorization":     {config.Resources.SessionTokens.AccessToken},
			"x-acf-sensor-data": {config.Resources.SensorData},
			"accept-language":   {"pl-PL"},
			"accept":            {"application/hal+json"},
			"content-type":      {"application/json;charset=UTF-8"},
			"x-api-key":         {"m79qyapn2kbucuv96ednvh22"},
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

		config.DefaultConfig.CreateEndTimestamp()

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
			return definederrors.ERROR_PLACEHOLDER
		}

		config.DefaultConfig.Log.Green("Order successfully placed!")

		if config.IfDiscounted {
			discountmenager.ConfirmUsage(config.DefaultConfig.TaskData.Discount, config.UsedDiscount)
		}

		body, err := httpclient.GetBodyString(res, &config.DefaultConfig.Log)

		if err != nil {
			continue
		}

		var orderResponse OrderResponse
		var paymentLink string

		json.Unmarshal([]byte(body), &orderResponse)

		if strings.EqualFold(config.DefaultConfig.Profile.Payment, "paypal") {
			paymentLink, err = GetPaypalLink(config, orderResponse.ThirdPartyPayment.Base64HTML)

			if err != nil {
				paymentLink = "https://www.youtube.com/watch?v=dQw4w9WgXcQ&t=1s&ab_channel=RickAstley"
			}
		}

		notifierPayload := webhookenginev2.NotifierPayload{
			BotFields: webhookenginev2.BotFields{
				PrivateWebhook: config.DefaultConfig.Profile.Webhook,
				PublicWebhook:  consts.PUBLIC_WEBHOOK,
				Img:            config.Cosmetics.Image,
				Title:          "Checkout successful!",
			},
			WebhookFields: webhookenginev2.WebhookFields{
				Name:        config.Cosmetics.Name,
				Site:        config.DefaultConfig.TaskData.Website,
				Region:      config.DefaultConfig.TaskData.Region,
				DicountCode: config.UsedDiscount,
				Speed:       config.DefaultConfig.GetCheckoutTime(),
				Discount:    config.UsedDiscount,
				Size:        config.Cosmetics.Size,
				Pid:         config.DefaultConfig.TaskData.Sku,
				Price:       config.Cosmetics.Price,
				Payment:     config.DefaultConfig.Profile.Payment,
				Mode:        config.DefaultConfig.TaskData.Mode,
				ProfileName: config.DefaultConfig.Profile.Name,
				TaskId:      config.DefaultConfig.TaskId,
				Email:       config.DefaultConfig.TaskData.Email,
				Password:    config.DefaultConfig.TaskData.Password,
				Proxy:       fmt.Sprintf("`%s`", config.DefaultConfig.Proxy.String),
				OrderNumber: fmt.Sprintf("[%s](%s)", GetOrderNumber(&orderResponse), paymentLink),
			},
		}

		successhandler.HandleSuccess(&notifierPayload)

		return nil
	}

	return config.DefaultConfig.Log.LogReturnErrorCustomText(definederrors.ERROR_TOO_MANY_RETRYS, definederrors.MESSAGE_TOO_MANY_RETRYS)
}
