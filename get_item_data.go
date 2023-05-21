package adidas_backend

import (
	"encoding/json"
	"fmt"
	"strings"
	httpclient "umbrella/internal/http_client"
	definederrors "umbrella/internal/utils/defined_errors"
	"umbrella/internal/utils/helpers"

	http "github.com/vimbing/fhttp"
)

func ItemData(config *Config) []Variant {
	state := config.TaskStates.Monitor

	config.DefaultConfig.Log.SetState(state.Name)

	for i := 0; i < state.Retry; i++ {
		// req, err := http.NewRequest("GET", fmt.Sprintf("https://api.3stripes.net/gw-api/v2/products/%s/availability?experiment_product_data=%s", config.DefaultConfig.TaskData.Sku, helpers.RandomString(8, true, false, true)), nil)
		req, err := http.NewRequest("GET", fmt.Sprintf("https://api.3stripes.net/gw-api/v2/products/%s/Availability/", strings.ToUpper(config.DefaultConfig.TaskData.Sku)), nil)

		if err != nil {
			config.DefaultConfig.Log.RedDelay(definederrors.MESSAGE_JSON_MARSHALING_ERROR)
			continue
		}

		req.Header = http.Header{
			"host":            {"api.3stripes.net"},
			"user-agent":      {fmt.Sprintf("adidas/2022.11.25.13.12CFNetwork/1390Darwin/22.0.0%s", helpers.RandomString(helpers.RandomInt(8, 25), true, false, true))},
			"x-market":        {"PL"},
			"accept-language": {fmt.Sprintf("pl-PL,pl;q=0.9,en-US;q=0.8,en;q=0.7,la;q=0.6,de;q=0.%s", helpers.RandomString(helpers.RandomInt(5, 35), true, false, true))},
			"accept":          {"application/hal+json"},
			"content-type":    {"application/json;charset=UTF-8"},
			"Cache-Control":   {"no-cache:max-age=0"},
			"Pragma":          {"no-cache"},
			"x-api-key":       {"m79qyapn2kbucuv96ednvh22"},
		}

		res, err := config.DefaultConfig.MonitorNetwork.Client.Do(req)

		if err != nil {
			config.DefaultConfig.Log.RedDelay(definederrors.MESSAGE_REQUEST_SENDING_ERROR)
			continue
		}

		defer res.Body.Close()

		if res.StatusCode != 200 {
			config.DefaultConfig.Log.StatusCodeErrorDelay(res.Status)
			return []Variant{}
		}

		body, err := httpclient.GetBodyString(res, &config.DefaultConfig.Log)

		if err != nil {
			continue
		}

		var getItemDataResponse GetItemDataResponse

		json.Unmarshal([]byte(body), &getItemDataResponse)

		var variants []Variant

		if len(getItemDataResponse.Embedded.Variations) < 1 {
			config.DefaultConfig.Log.YellowDelay("No variants for item, maybe not loaded...")
			return variants
		}

		for _, v := range getItemDataResponse.Embedded.Variations {
			if v.Orderable {
				variants = append(variants, Variant{
					SizePid: v.VariationProductID,
					Size:    v.Size,
					Pid:     getItemDataResponse.ProductID,
				})
			}
		}

		return variants
	}

	return []Variant{}
}
