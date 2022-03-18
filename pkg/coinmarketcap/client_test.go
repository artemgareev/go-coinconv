package coinmarketcap

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"go-coinconv/pkg/coinmarketcap/protocol"
)

type httpClientMock struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (h *httpClientMock) Do(req *http.Request) (*http.Response, error) {
	return h.DoFunc(req)
}

var (
	priceConversionSuccessResponse = []byte(`{
    "status": {
        "timestamp": "2022-03-17T17:22:16.708Z",
        "error_code": 0,
        "error_message": null,
        "elapsed": 35,
        "credit_count": 1,
        "notice": null
    },
    "data": [{
        "id": 1,
        "symbol": "BTC",
        "name": "Bitcoin",
        "amount": 1,
        "last_updated": "2022-03-17T17:22:00.000Z",
        "quote": {
            "ETH": {
                "price": 14.59971897907295,
                "last_updated": "2022-03-17T17:22:00.000Z"
            }
        }
    }]
}`)

	priceConversionBadResponse = []byte(`{
	"status": {
		"timestamp": "2018-06-02T22:51:28.209Z",
		"error_code": 500,
		"error_message": "An internal server error occurred",
		"elapsed": 10,
		"credit_count": 0
	}
}`)
	successResponseStruct protocol.PriceConversionResponse
	badResponseStruct     protocol.PriceConversionResponse
)

func init() {
	if err := json.Unmarshal(priceConversionSuccessResponse, &successResponseStruct); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(priceConversionBadResponse, &badResponseStruct); err != nil {
		panic(err)
	}
}

func TestClient_PriceConversion1(t *testing.T) {
	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	type fields struct {
		httpClient httpClient
	}
	type args struct {
		ctx     context.Context
		request *protocol.PriceConversionQueryParameters
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *protocol.PriceConversionResponse
		wantErr bool
	}{
		{
			name: "success_response",
			fields: fields{
				httpClient: &httpClientMock{func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader(priceConversionSuccessResponse)),
					}, nil
				}},
			},
			args: args{
				ctx:     context.Background(),
				request: &protocol.PriceConversionQueryParameters{},
			},
			want:    &successResponseStruct,
			wantErr: false,
		},
		{
			name: "bad_response",
			fields: fields{
				httpClient: &httpClientMock{func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: 500,
						Body:       ioutil.NopCloser(bytes.NewReader(priceConversionBadResponse)),
					}, nil
				}},
			},
			args: args{
				ctx:     context.Background(),
				request: &protocol.PriceConversionQueryParameters{},
			},
			want:    &badResponseStruct,
			wantErr: false,
		},
		{
			name: "malformed_response",
			fields: fields{
				httpClient: &httpClientMock{func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`something went wrong`))),
					}, nil
				}},
			},
			args: args{
				ctx:     context.Background(),
				request: &protocol.PriceConversionQueryParameters{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "unknown_status_code_response",
			fields: fields{
				httpClient: &httpClientMock{func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: 501,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(`something went wrong`))),
					}, nil
				}},
			},
			args: args{
				ctx:     context.Background(),
				request: &protocol.PriceConversionQueryParameters{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ctx_done",
			fields: fields{
				httpClient: &http.Client{},
			},
			args: args{
				ctx:     canceledCtx,
				request: &protocol.PriceConversionQueryParameters{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				httpClient: tt.fields.httpClient,
				baseURL:    "https://google.com",
				apiKey:     "token",
			}
			got, err := c.PriceConversion(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("PriceConversion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PriceConversion() got = %v, want %v", got, tt.want)
			}
		})
	}
}
