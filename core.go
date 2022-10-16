package main

import (
	"encoding/json"
	"errors"

	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Client .
type Client struct {
	apiKey     string
	baseURL    string
	HTTPClient *http.Client
}

// NewClient creates new client with given API key
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
		baseURL: "https://api.facest.io/v1",
	}
}

// Rectangle .
type Rectangle struct {
	Top    int `json:"top"`
	Left   int `json:"left"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

// Content-type and body should be already added to req
func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

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

// GetFaces returns trained faces and their images
func (c *Client) GetFaces(options *FacesListOptions) (*FacesList, error) {
	limit := 100
	page := 1
	if options != nil {
		limit = options.Limit
		page = options.Page
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/faces?limit=%d&page=%d", c.baseURL, limit, page), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	res := FacesList{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetFace returns face object and its images by face_token
func (c *Client) GetFace(faceToken string) (*Face, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/faces/%s", c.baseURL, url.PathEscape(faceToken)), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	res := Face{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
