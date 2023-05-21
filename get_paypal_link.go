package adidas_backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	httpclient "umbrella/internal/http_client"
	definederrors "umbrella/internal/utils/defined_errors"

	http "github.com/vimbing/fhttp"
)

func GetPaypalLink(config *Config, adidasBody string) (string, error) {
	state := config.TaskStates.Payment

	config.DefaultConfig.Log.SetState(state.Name)

	for i := 0; i < state.Retry; i++ {
		addPaypalLinkToApiPayload := AddPaypalLinkToApiPayload{
			Body: adidasBody,
		}

		payload, err := json.Marshal(addPaypalLinkToApiPayload)

		if err != nil {
			config.DefaultConfig.Log.RedDelay(definederrors.MESSAGE_JSON_MARSHALING_ERROR)
			continue
		}

		req, err := http.NewRequest("POST", "https://api.umbrellaio.dev/adidas/paypal/add", bytes.NewBuffer(payload))

		if err != nil {
			config.DefaultConfig.Log.RedDelay(definederrors.MESSAGE_JSON_MARSHALING_ERROR)
			continue
		}

		req.Header = http.Header{
			"x-umbrella-auth-token": {"4dbffc9823cc1f78d620f1b8af50a74c2445a0821de72bf9eb1a1d5cf479e5e1"},
			"content-type":          {"application/json"},
		}

		res, err := config.DefaultConfig.Network.Client.Do(req)

		if err != nil {
			config.DefaultConfig.Log.RedDelay(definederrors.MESSAGE_REQUEST_SENDING_ERROR)
			continue
		}

		defer res.Body.Close()

		if res.StatusCode != 200 {
			config.DefaultConfig.Log.StatusCodeErrorDelay(res.Status)
			return "", definederrors.ERROR_PLACEHOLDER
		}

		body, err := httpclient.GetBodyString(res, &config.DefaultConfig.Log)

		if err != nil {
			continue
		}

		var umbrellaApiLinkAddResponse UmbrellaApiLinkAddResponse

		json.Unmarshal([]byte(body), &umbrellaApiLinkAddResponse)

		return fmt.Sprintf("https://api.umbrellaio.dev/adidas/paypal/get/%s", umbrellaApiLinkAddResponse.URL), nil
	}

	return "", definederrors.ERROR_PLACEHOLDER
}
