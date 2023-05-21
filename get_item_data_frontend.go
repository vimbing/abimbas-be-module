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

func ItemDataFrontend(config *Config) []Variant {
	state := config.TaskStates.Monitor

	config.DefaultConfig.Log.SetState(state.Name)

	for i := 0; i < state.Retry; i++ {
		// req, err := http.NewRequest("GET", fmt.Sprintf("https://www.adidas.pl/api/products/%s/availability", strings.ToUpper(config.DefaultConfig.TaskData.Sku)), nil)
		req, err := http.NewRequest("GET", fmt.Sprintf("https://www.adidas.pl/on/demandware.store/Sites-adidas-PL-Site/pl_PL/Product-GetProductAvailabilityJSON?pid=%s&s=%s", strings.ToUpper(config.DefaultConfig.TaskData.Sku), helpers.RandomString(helpers.RandomInt(0, 25), true, false, true)), nil)

		if err != nil {
			config.DefaultConfig.Log.RedDelay(definederrors.MESSAGE_JSON_MARSHALING_ERROR)
			continue
		}

		req.Header = http.Header{
			"authority":          {"www.adidas.pl"},
			"accept":             {"*/*"},
			"accept-language":    {"pl-PL,pl;q=0.9,en-US;q=0.8,en;q=0.7,la;q=0.6,de;q=0.5"},
			"content-type":       {"application/json"},
			"dnt":                {"1"},
			"glassversion":       {"c45f0f7"},
			"sec-ch-ua":          {"\"Not?A_Brand\";v=\"8\",\"Chromium\";v=\"108\",\"GoogleChrome\";v=\"108\""},
			"sec-ch-ua-mobile":   {"?0"},
			"sec-ch-ua-platform": {"\"Windows\""},
			"sec-fetch-dest":     {"empty"},
			"sec-fetch-mode":     {"cors"},
			"sec-fetch-site":     {"same-origin"},
			"user-agent":         {"Mozilla/5.0(WindowsNT10.0;Win64;x64)AppleWebKit/537.36(KHTML,likeGecko)Chrome/108.0.0.0Safari/537.36"},
			"x-instana-l":        {"1,correlationType=web;correlationId=9a5890e22892fefd"},
			"x-instana-s":        {"9a5890e22892fefd"},
			"x-instana-t":        {"9a5890e22892fefd"},
			http.HeaderOrderKey: {
				"authority",
				"accept",
				"accept-language",
				"content-type",
				"dnt",
				"glassversion",
				"referer",
				"sec-ch-ua",
				"sec-ch-ua-mobile",
				"sec-ch-ua-platform",
				"sec-fetch-dest",
				"sec-fetch-mode",
				"sec-fetch-site",
				"user-agent",
				"x-instana-l",
				"x-instana-s",
				"x-instana-t",
			},
			http.PHeaderOrderKey: {
				":method",
				":authority",
				":scheme",
				":path",
			},
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

		if res.StatusCode == 400 {
			return []Variant{}
		}

		body, err := httpclient.GetBodyString(res, &config.DefaultConfig.Log)

		if err != nil {
			continue
		}

		var itemDataFrontendResponse ItemDataFrontendResponse

		json.Unmarshal([]byte(body), &itemDataFrontendResponse)

		var variants []Variant

		if len(itemDataFrontendResponse.VariationList) < 1 {
			config.DefaultConfig.Log.YellowDelay("No variants for item, maybe not loaded...")
			return variants
		}

		for _, v := range itemDataFrontendResponse.VariationList {
			if v.Availability > 0 {
				variants = append(variants, Variant{
					SizePid: v.Sku,
					Size:    v.Size,
					Pid:     itemDataFrontendResponse.ID,
				})
			}
		}

		return variants
	}

	return []Variant{}
}
