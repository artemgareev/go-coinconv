package coinmarketcap

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/go-querystring/query"

	"go-coinconv/pkg/coinmarketcap/protocol"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	httpClient httpClient
	baseURL    string
	apiKey     string
}

func NewClient(apiKey, baseAPIUrl string) *Client {
	return &Client{
		httpClient: &http.Client{},
		baseURL:    baseAPIUrl,
		apiKey:     apiKey,
	}
}

// PriceConversion https://coinmarketcap.com/api/v1/#operation/getV2ToolsPriceconversion
func (c *Client) PriceConversion(ctx context.Context, request *protocol.PriceConversionQueryParameters) (*protocol.PriceConversionResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/v2/tools/price-conversion", nil)
	if err != nil {
		return nil, err
	}
	queryVals, err := query.Values(request)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = queryVals.Encode()
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CMC_PRO_API_KEY", c.apiKey)

	var priceConversionResponse protocol.PriceConversionResponse
	if err := c.doRequest(ctx, req, &priceConversionResponse); err != nil {
		return nil, err
	}

	return &priceConversionResponse, nil
}

func (c *Client) doRequest(ctx context.Context, req *http.Request, v interface{}) error {
	req = req.WithContext(ctx)

	httpResp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode >= 200 && httpResp.StatusCode <= 500 {
		return json.NewDecoder(httpResp.Body).Decode(v)
	}

	byteResponse, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}

	return fmt.Errorf("%d bad response: %s", int32(httpResp.StatusCode), string(byteResponse))
}
