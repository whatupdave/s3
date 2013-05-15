package s3util

import (
	"net/http"
	"time"
)

// Deletes object at url. An HTTP status other than 200 is
// considered an error.
//
// If c is nil, Open uses DefaultConfig.
func Delete(url string, c *Config) error {
	if c == nil {
		c = DefaultConfig
	}

	r, _ := http.NewRequest("DELETE", url, nil)
	r.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	c.Sign(r, *c.Keys)
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		return newRespError(resp)
	}
	resp.Body.Close()
	return nil
}
