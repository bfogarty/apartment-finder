package main

// ShortenReq is the request to the v4 /shorten endpoint
type ShortenReq struct {
	LongURL string `json:"long_url"`
}

// ShortenResp is the response from the v4 /shorten endpoint
type ShortenResp struct {
	Link string `json:"link"`
}
