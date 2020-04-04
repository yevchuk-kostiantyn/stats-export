package edbo

import (
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Client struct {
	baseURL string
	client  *http.Client
}

type config struct {
	HttpTimeout time.Duration `split_words:"true" required:"true" default:"20s"`
	EDBOUrl     string        `split_words:"true" required:"true"`
}

func NewClient() (*Client, error) {
	var c config
	if err := envconfig.Process("", &c); err != nil {
		return nil, err
	}
	return &Client{
		baseURL: c.EDBOUrl,
		client:  &http.Client{Timeout: c.HttpTimeout},
	}, nil
}
