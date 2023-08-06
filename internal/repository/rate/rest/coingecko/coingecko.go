package coingecko

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"gses2-app/internal/core/port"
	"gses2-app/internal/repository/rate/rest"
)

// Represents data type for JSON response
type Response struct {
	Bitcoin struct {
		UAH float64 `json:"uah"`
	} `json:"bitcoin"`
}

var (
	ErrHTTPRequestFailure       = errors.New("http request failure")
	ErrUnexpectedStatusCode     = errors.New("unexpected status code")
	ErrUnexpectedResponseFormat = errors.New("unexpected response format")
)

const (
	_providerName = "CoingeckoRateProvider"
)

type CoingeckoAPIConfig struct {
	URL string `default:"https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=uah"`
}

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type CoingeckoProvider struct {
	config CoingeckoAPIConfig
}

func NewProvider(
	logger port.Logger,
	config CoingeckoAPIConfig,
	httpClient HTTPClient,
) *rest.AbstractProvider {
	return rest.NewProvider(
		logger,
		&CoingeckoProvider{
			config: config,
		},
		httpClient,
	)
}

func (p *CoingeckoProvider) URL() string {
	return p.config.URL
}

func (p *CoingeckoProvider) Name() string {
	return _providerName
}

func (p *CoingeckoProvider) ExtractRate(resp *http.Response) (port.Rate, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var data Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, errors.Join(err, ErrUnexpectedResponseFormat)
	}

	return port.Rate(data.Bitcoin.UAH), nil
}
