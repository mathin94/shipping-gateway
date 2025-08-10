package biteship

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	logger     *logrus.Logger
}

const (
	V1GetArea           = "/v1/maps/areas"
	V1GetCourier        = "/v1/couriers"
	V1GetCourierRates   = "/v1/rates/couriers"
	V1TrackingByID      = "/v1/tracking/%s"
	V1TrackingByWaybill = "/v1/tracking/%s/couriers/%s"
)

func NewClient(cfg *viper.Viper, logger *logrus.Logger) *Client {
	return &Client{
		baseURL: cfg.GetString("biteship.base_url"),
		apiKey:  cfg.GetString("biteship.api_key"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second, // Set a timeout for HTTP requests
		},
		logger: logger,
	}
}

func (c *Client) GetRequest(endpoint string, queryParams map[string]string) (int, []byte, error) {
	bRes := make([]byte, 0)

	req, err := http.NewRequest(http.MethodGet, c.baseURL+endpoint, nil)
	if err != nil {
		return http.StatusInternalServerError, bRes, err
	}

	// Add query parameters
	q := req.URL.Query()
	for k, v := range queryParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	// Add API key to header
	req.Header.Set("Authorization", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Errorf("Error making GET request to %s: %v", req.URL.String(), err)
		return http.StatusInternalServerError, bRes, err
	}
	defer resp.Body.Close()
	bRes, _ = io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		// log http request response with all data including headers and body
		c.logger.WithField("request", queryParams).
			WithField("response", bRes).
			WithField("headers", req.Header).
			Errorf("Unexpected status code: %d", resp.StatusCode)

		return resp.StatusCode, bRes, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp.StatusCode, bRes, nil
}

func (c *Client) PostRequest(endpoint string, body any) (int, []byte, error) {
	reqBody, err := json.Marshal(body)
	bRes := make([]byte, 0)
	if err != nil {
		c.logger.Errorf("Error marshalling request body: %v", err)
		return http.StatusInternalServerError, bRes, err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+endpoint, bytes.NewReader(reqBody))
	if err != nil {
		c.logger.Errorf("Error creating POST request: %v", err)
		return http.StatusInternalServerError, bRes, err
	}

	req.Header.Set("Authorization", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Errorf("Error making POST request to %s: %v", req.URL.String(), err)
		return http.StatusInternalServerError, bRes, err
	}
	defer resp.Body.Close()

	bRes, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		// log http request response with all data including headers and body
		c.logger.WithField("request", string(reqBody)).
			WithField("response", bRes).
			WithField("headers", req.Header).
			Errorf("Unexpected status code: %d", resp.StatusCode)

		return resp.StatusCode, bRes, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp.StatusCode, bRes, nil
}

func (c *Client) SearchAreas(query string) (resp *AreaResponse, err *ErrorResponse) {
	queryParams := map[string]string{
		"countries": "ID",
		"input":     query,
		"type":      "single",
	}

	_, bRes, errResp := c.GetRequest(V1GetArea, queryParams)
	if errResp != nil {
		c.logger.Errorf("Error getting area: %v, biteship resp : %s", err, bRes)
		return nil, ErrorResponseFromBytes(bRes)
	}

	if err := json.Unmarshal(bRes, &resp); err != nil {
		c.logger.Errorf("Error unmarshalling response: %s, error : %v", bRes, err)
		return nil, NewErrorResponse(ErrInvalidParsingResponse, err.Error())
	}

	return resp, nil
}

func (c *Client) GetCourierRates(request RateRequest) (resp *RateResponse, errResp *ErrorResponse) {
	// convert request to json string for logging
	requestJSON, _ := json.Marshal(request)
	c.logger.Debugf("Requesting courier rates with request: %s", requestJSON)
	_, bRes, err := c.PostRequest(V1GetCourierRates, request)
	if err != nil {
		c.logger.Errorf("Error getting courier rates: %v, biteship resp : %s", err, bRes)
		return nil, ErrorResponseFromBytes(bRes)
	}

	if err := json.Unmarshal(bRes, &resp); err != nil {
		c.logger.Errorf("Error unmarshalling response: %s, error : %v", bRes, err)
		return nil, NewErrorResponse(ErrInvalidParsingResponse, err.Error())
	}

	return resp, nil
}

func (c *Client) GetTrackingByWaybill(waybill, courier string) (resp *TrackingResponse, errResp *ErrorResponse) {
	c.logger.Debugf("Requesting tracking by waybill: %s, courier: %s", waybill, courier)
	endpoint := fmt.Sprintf(V1TrackingByWaybill, waybill, courier)
	_, bRes, err := c.GetRequest(endpoint, nil)
	if err != nil {
		c.logger.Errorf("Error getting tracking by waybill: %v, biteship resp : %s", err, bRes)
		return nil, ErrorResponseFromBytes(bRes)
	}

	if err := json.Unmarshal(bRes, &resp); err != nil {
		c.logger.Errorf("Error unmarshalling response: %s, error : %v", bRes, err)
		return nil, NewErrorResponse(ErrInvalidParsingResponse, err.Error())
	}

	return resp, nil
}
