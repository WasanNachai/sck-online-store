package point

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type PointGateway struct {
	PointEndpoint string
	HTTPClient    *http.Client
}

func (gateway PointGateway) getHTTPClient() *http.Client {
	if gateway.HTTPClient != nil {
		return gateway.HTTPClient
	}

	return &http.Client{
		Timeout: 3 * time.Second,
	}
}

func (gateway PointGateway) buildEndpoint(path string) string {
	return strings.TrimRight(gateway.PointEndpoint, "/") + path
}

func (gateway PointGateway) GetPoints(
	ctx context.Context,
	uid int,
) ([]Point, error) {
	endpoint := gateway.buildEndpoint("/api/v1/point")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint,
		nil,
	)
	if err != nil {
		return []Point{}, fmt.Errorf(
			"create get-points request: %w",
			err,
		)
	}

	response, err := gateway.getHTTPClient().Do(req)
	if err != nil {
		return []Point{}, fmt.Errorf(
			"call point-service get-points endpoint: %w",
			err,
		)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(response.Body)

		return []Point{}, fmt.Errorf(
			"point-service get-points returned status %d: %s",
			response.StatusCode,
			string(responseBody),
		)
	}

	var gatewayResponse []Point

	if err := json.NewDecoder(response.Body).Decode(
		&gatewayResponse,
	); err != nil {
		return []Point{}, fmt.Errorf(
			"decode get-points response: %w",
			err,
		)
	}

	return gatewayResponse, nil
}

func (gateway PointGateway) CreatePoint(
	ctx context.Context,
	uid int,
	body Point,
) (Point, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return Point{}, fmt.Errorf(
			"marshal create-point request: %w",
			err,
		)
	}

	endpoint := gateway.buildEndpoint("/api/v1/point")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		bytes.NewReader(data),
	)
	if err != nil {
		return Point{}, fmt.Errorf(
			"create point-service request: %w",
			err,
		)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	response, err := gateway.getHTTPClient().Do(req)
	if err != nil {
		return Point{}, fmt.Errorf(
			"call point-service create-point endpoint: %w",
			err,
		)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK &&
		response.StatusCode != http.StatusCreated {
		responseBody, _ := io.ReadAll(response.Body)

		return Point{}, fmt.Errorf(
			"point-service create-point returned status %d: %s",
			response.StatusCode,
			string(responseBody),
		)
	}

	var gatewayResponse Point

	if err := json.NewDecoder(response.Body).Decode(
		&gatewayResponse,
	); err != nil {
		return Point{}, fmt.Errorf(
			"decode create-point response: %w",
			err,
		)
	}

	return gatewayResponse, nil
}

func (gateway PointGateway) CalculateEarnedPoints(
	ctx context.Context,
	amountTHB float64,
) (CalculatePointResponse, error) {
	if amountTHB <= 0 {
		return CalculatePointResponse{
			EarnedPoints:    0,
			AmountTHB:       amountTHB,
			RateTHBPerPoint: 50,
		}, nil
	}

	if gateway.PointEndpoint == "" {
		return CalculatePointResponse{},
			fmt.Errorf("point-service endpoint is empty")
	}

	requestBody := CalculatePointRequest{
		AmountTHB: amountTHB,
	}

	data, err := json.Marshal(requestBody)
	if err != nil {
		return CalculatePointResponse{}, fmt.Errorf(
			"marshal calculate-point request: %w",
			err,
		)
	}

	endpoint := gateway.buildEndpoint(
		"/api/v1/point/calculate",
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		bytes.NewReader(data),
	)
	if err != nil {
		return CalculatePointResponse{}, fmt.Errorf(
			"create calculate-point request: %w",
			err,
		)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	response, err := gateway.getHTTPClient().Do(req)
	if err != nil {
		return CalculatePointResponse{}, fmt.Errorf(
			"call point-service calculate endpoint: %w",
			err,
		)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK &&
		response.StatusCode != http.StatusCreated {
		responseBody, _ := io.ReadAll(response.Body)

		return CalculatePointResponse{}, fmt.Errorf(
			"point-service calculate endpoint returned status %d: %s",
			response.StatusCode,
			string(responseBody),
		)
	}

	var gatewayResponse CalculatePointResponse

	if err := json.NewDecoder(response.Body).Decode(
		&gatewayResponse,
	); err != nil {
		return CalculatePointResponse{}, fmt.Errorf(
			"decode calculate-point response: %w",
			err,
		)
	}

	if gatewayResponse.EarnedPoints < 0 {
		return CalculatePointResponse{}, fmt.Errorf(
			"point-service returned invalid earned points: %d",
			gatewayResponse.EarnedPoints,
		)
	}

	return gatewayResponse, nil
}
