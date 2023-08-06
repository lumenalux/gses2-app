package binance

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"gses2-app/internal/core/port"
	"gses2-app/internal/repository/rate/rest"
)

var (
	ErrHTTPRequestFailure           = errors.New("http request failure")
	ErrUnexpectedStatusCode         = errors.New("unexpected status code")
	ErrUnexpectedResponseFormat     = errors.New("unexpected response format")
	ErrUnexpectedExchangeRateFormat = errors.New("unexpected exchange rate format")
)

const (
	_providerName     = "BinanceRateProvider"
	_firstItemIndex   = 0
	_minResponseItems = 5
	_rateIndex        = 4
)

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type BinanceAPIConfig struct {
	URL string `default:"https://api.binance.com/api/v3/klines?symbol=BTCUAH&interval=1s&limit=1"`
}

type BinanceProvider struct {
	config BinanceAPIConfig
}

func NewProvider(
	logger port.Logger,
	config BinanceAPIConfig,
	httpClient HTTPClient,
) *rest.AbstractProvider {
	return rest.NewProvider(
		logger,
		&BinanceProvider{
			config: config,
		},
		httpClient,
	)
}

func (p *BinanceProvider) URL() string {
	return p.config.URL
}

func (p *BinanceProvider) Name() string {
	return _providerName
}

func (p *BinanceProvider) ExtractRate(resp *http.Response) (port.Rate, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var data [][]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, err
	}

	if len(data) == 0 || len(data[_firstItemIndex]) < _minResponseItems {
		return 0, ErrUnexpectedResponseFormat
	}

	exchangeRate, ok := data[_firstItemIndex][_rateIndex].(string)
	if !ok {
		return 0, ErrUnexpectedExchangeRateFormat
	}

	rateValue, err := strconv.ParseFloat(exchangeRate, 64)
	if err != nil {
		return 0, errors.Join(err, ErrUnexpectedExchangeRateFormat)
	}

	return port.Rate(rateValue), nil
}
