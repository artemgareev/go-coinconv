package provider

import (
	"context"
	"errors"
	"strings"

	"go-coinconv/pkg/coinmarketcap"
	"go-coinconv/pkg/coinmarketcap/protocol"
)

type coinMarketCapProvider struct {
	client *coinmarketcap.Client
}

func NewCoinMarketCapProvider(cmcClient *coinmarketcap.Client) Convertor {
	return &coinMarketCapProvider{
		client: cmcClient,
	}
}

func (c *coinMarketCapProvider) Convert(ctx context.Context, amount float64, symbolFrom, symbolTo string) (float64, error) {
	request := &protocol.PriceConversionQueryParameters{
		Amount:  amount,
		Symbol:  &symbolFrom,
		Convert: &symbolTo,
	}

	response, err := c.client.PriceConversion(ctx, request)
	// todo: consider logging, metrics
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, errors.New("request timeout")
		}
		return 0, errors.New("something went wrong")
	}

	if response.Status.ErrorCode != 0 {
		return 0, errors.New(*response.Status.ErrorMessage)
	}

	if len(response.Data) < 1 {
		return 0, errors.New("bad response")
	}
	if quote, ok := response.Data[0].Quote[strings.ToUpper(symbolTo)]; ok {
		return quote.Price, nil
	}

	return 0, errors.New("bad response")
}
