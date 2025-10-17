package weatherapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Client struct {
	http    *http.Client
	apiKey  string
	lang    string
	timeout time.Duration
}

func NewClient(apikey, lang string, timeout time.Duration) *Client {
	return &Client{
		http:    &http.Client{Timeout: timeout},
		apiKey:  apikey,
		lang:    lang,
		timeout: timeout,
	}
}

func (c *Client) Forecast(ctx context.Context, query string, days int, aqi, alerts bool) (*Weather, error) {
	u, _ := url.Parse(os.Getenv("BASE_URL"))
	q := u.Query()
	q.Set("key", c.apiKey)
	q.Set("q", query)
	q.Set("days", strconv.Itoa(days))
	q.Set("lang", c.lang)
	q.Set("aqi", boolToYesNo(aqi))
	q.Set("alerts", boolToYesNo(alerts))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	// Uncomment next line in case the provider requires acustom User-Agent
	// req.Header.Set("User-Agent", "weather-cli/1.0 (+github.com/username/cliweather)")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("weatherapi: http %d", resp.StatusCode)
	}

	var w Weather
	if err := json.NewDecoder(resp.Body).Decode(&w); err != nil {
		return nil, err
	}
	return &w, nil
}

func boolToYesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}
