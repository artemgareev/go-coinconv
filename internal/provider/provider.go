package provider

import "context"

type Convertor interface {
	Convert(ctx context.Context, amount float64, symbolFrom, symbolTo string) (float64, error)
}
