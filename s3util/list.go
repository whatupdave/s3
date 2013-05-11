package s3util

import (
	"encoding/xml"
	"net/http"
	"strconv"
	"time"
)

type Owner struct {
	ID          string
	DisplayName string
}

type Key struct {
	Key          string
	LastModified string
	Size         int64
	// ETag gives the hex-encoded MD5 sum of the contents,
	// surrounded with double-quotes.
	ETag         string
	StorageClass string
	Owner        Owner
}

type ListResp struct {
	Name      string
	Prefix    string
	Delimiter string
	Marker    string
	MaxKeys   int
	// IsTruncated is true if the results have been truncated because
	// there are more keys and prefixes than can fit in MaxKeys.
	// N.B. this is the opposite sense to that documented (incorrectly) in
	// http://goo.gl/YjQTc
	IsTruncated    bool
	Contents       []Key
	CommonPrefixes []string `xml:">Prefix"`
}

func List(url, prefix, marker string, max int, c *Config) (*ListResp, error) {
	if c == nil {
		c = DefaultConfig
	}

	r, _ := http.NewRequest("GET", url+"/?prefix="+prefix+"&marker="+marker+"&max-keys="+strconv.Itoa(max), nil)
	r.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	c.Sign(r, *c.Keys)
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, newRespError(resp)
	}
	lr := &ListResp{}
	err = xml.NewDecoder(resp.Body).Decode(lr)
	if err != nil {
		return nil, err
	}
	return lr, nil
}
