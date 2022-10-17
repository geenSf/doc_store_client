package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"

	"fmt"
	"net/http"
	"os"
	"time"
)

// Client .
type Client struct {
	HTTPClient *http.Client
	baseURL    string
	Transport  *http.Transport
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

// creates new client with given API key
func NewClient(baseURL string, certPool *x509.CertPool) *Client {
	return &Client{

		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
		baseURL:   "https://" + baseURL,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: certPool}},
	}

}

// implements GET request
func (c *Client) PostReq(filename string, user string, pass string) (*http.Response, error) {

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	reader := bufio.NewReader(f)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.baseURL, filename), reader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.SetBasicAuth(user, pass)

	var res *http.Response

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// implements GET request
func (c *Client) GetReq(key string) (*http.Response, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.baseURL, key), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json; charset=utf-8")

	var res *http.Response

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// Content-type and body should be already added to req
func (c *Client) sendRequest(req *http.Request, v interface{}) error {

	//req.Header.Set("Accept", "application/json; charset=utf-8")
	//req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// Try to unmarshall into errorResponse
	if res.StatusCode != http.StatusOK {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	// Unmarshall and populate v
	fullResponse := successResponse{
		Data: v,
	}

	if err = json.NewDecoder(res.Body).Decode(&fullResponse); err != nil {
		return err
	}

	return nil
}
