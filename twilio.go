package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var messagesEndpoint = "https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json"

// TwilioClient is a client for the Twilio API.
type TwilioClient struct {
	accountSID string
	authToken  string
	client     *http.Client
}

// NewTwilioClient returns a new Twilio API client.
func NewTwilioClient(accountSID string, authToken string) *TwilioClient {
	return &TwilioClient{
		accountSID: accountSID,
		authToken:  authToken,
		client:     &http.Client{},
	}
}

// SendSMS sends an SMS using the Twilio API.
func (t *TwilioClient) SendSMS(to string, from string, body string) error {
	msg := url.Values{"To": {to}, "From": {from}, "Body": {body}}
	msgStr := strings.NewReader(msg.Encode())

	endpoint := fmt.Sprintf(messagesEndpoint, t.accountSID)

	req, err := http.NewRequest("POST", endpoint, msgStr)
	if err != nil {
		return ErrSendingSMS
	}

	req.SetBasicAuth(t.accountSID, t.authToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	_, err = t.client.Do(req)

	if err != nil {
		return ErrSendingSMS
	}

	log.Printf("Sent SMS to %s", to)

	return nil
}
