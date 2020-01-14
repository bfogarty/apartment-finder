package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var shortenEndpoint = "https://api-ssl.bitly.com/v4/shorten"

// BitlyClient is a client for the Bitly API.
type BitlyClient struct {
	accessToken string
	client      *http.Client
}

// ShortenReq is the request to the v4 /shorten endpoint
type ShortenReq struct {
	LongURL string `json:"long_url"`
}

// ShortenResp is the response from the v4 /shorten endpoint
type ShortenResp struct {
	Link string `json:"link"`
}

// NewBitlyClient returns a new Bitly API client.
func NewBitlyClient(accessToken string) *BitlyClient {
	return &BitlyClient{
		accessToken: accessToken,
		client:      &http.Client{},
	}
}

// ShortenLink shortens a long URL using the Bitly API.
func (b *BitlyClient) ShortenLink(link string) (string, error) {
	body, err := json.Marshal(&ShortenReq{LongURL: link})
	if err != nil {
		return "", ErrShorteningLink
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", shortenEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return "", ErrShorteningLink
	}

	authToken := fmt.Sprintf("Bearer %s", b.accessToken)
	req.Header.Add("Authorization", authToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", ErrShorteningLink
	}

	data := &ShortenResp{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(data)
	if err != nil {
		return "", ErrShorteningLink
	}

	return data.Link, nil
}
