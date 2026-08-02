package point

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	params := url.Values{}
	params.Set("userId", fmt.Sprintf("%d", uid))
	endpoint = endpoint + "?" + params.Encode()

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

func (gateway PointGateway) GetPointSummary(
	ctx context.Context,
	uid int,
) (PointServiceSummary, error) {
	endpoint := gateway.buildEndpoint(
		fmt.Sprintf("/api/v1/point/summary/%d", uid),
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint,
		nil,
	)
	if err != nil {
		return PointServiceSummary{}, fmt.Errorf(
			"create get-point-summary request: %w",
			err,
		)
	}

	response, err := gateway.getHTTPClient().Do(req)
	if err != nil {
		return PointServiceSummary{}, fmt.Errorf(
			"call point-service get-point-summary endpoint: %w",
			err,
		)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(response.Body)
		return PointServiceSummary{}, fmt.Errorf(
			"point-service get-point-summary returned status %d: %s",
			response.StatusCode,
			string(responseBody),
		)
	}

	var gatewayResponse PointServiceSummary
	if err := json.NewDecoder(response.Body).Decode(
		&gatewayResponse,
	); err != nil {
		return PointServiceSummary{}, fmt.Errorf(
			"decode get-point-summary response: %w",
			err,
		)
	}

	return gatewayResponse, nil
}

func (gateway PointGateway) CreatePendingEarnPoint(
	ctx context.Context,
	uid int,
	body CreatePendingPointRequest,
) (Point, error) {
	_ = uid

	data, err := json.Marshal(body)
	if err != nil {
		return Point{}, fmt.Errorf(
			"marshal create-pending-point request: %w",
			err,
		)
	}

	endpoint := gateway.buildEndpoint("/api/v1/point/earn-pending")
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		bytes.NewReader(data),
	)
	if err != nil {
		return Point{}, fmt.Errorf(
			"create pending-point request: %w",
			err,
		)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	response, err := gateway.getHTTPClient().Do(req)
	if err != nil {
		return Point{}, fmt.Errorf(
			"call point-service create-pending-point endpoint: %w",
			err,
		)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK &&
		response.StatusCode != http.StatusCreated {
		responseBody, _ := io.ReadAll(response.Body)
		return Point{}, fmt.Errorf(
			"point-service create-pending-point returned status %d: %s",
			response.StatusCode,
			string(responseBody),
		)
	}

	var gatewayResponse Point
	if err := json.NewDecoder(response.Body).Decode(
		&gatewayResponse,
	); err != nil {
		return Point{}, fmt.Errorf(
			"decode create-pending-point response: %w",
			err,
		)
	}

	return gatewayResponse, nil
}

func (gateway PointGateway) ApprovePoint(
	ctx context.Context,
	pointID int,
	body PointTransitionRequest,
) (Point, error) {
	return gateway.transitionPoint(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("/api/v1/point/%d/approve", pointID),
		body,
		"approve-point",
	)
}

func (gateway PointGateway) RedeemPoint(
	ctx context.Context,
	pointID int,
	body PointTransitionRequest,
) (Point, error) {
	return gateway.transitionPoint(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("/api/v1/point/%d/redeem", pointID),
		body,
		"redeem-point",
	)
}

func (gateway PointGateway) ExpirePoints(
	ctx context.Context,
	body ExpirePointRequest,
) (ExpirePointResponse, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return ExpirePointResponse{}, fmt.Errorf(
			"marshal expire-point request: %w",
			err,
		)
	}

	endpoint := gateway.buildEndpoint("/api/v1/point/expire")
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		endpoint,
		bytes.NewReader(data),
	)
	if err != nil {
		return ExpirePointResponse{}, fmt.Errorf(
			"create expire-point request: %w",
			err,
		)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	response, err := gateway.getHTTPClient().Do(req)
	if err != nil {
		return ExpirePointResponse{}, fmt.Errorf(
			"call point-service expire-point endpoint: %w",
			err,
		)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(response.Body)
		return ExpirePointResponse{}, fmt.Errorf(
			"point-service expire-point returned status %d: %s",
			response.StatusCode,
			string(responseBody),
		)
	}

	var gatewayResponse ExpirePointResponse
	if err := json.NewDecoder(response.Body).Decode(
		&gatewayResponse,
	); err != nil {
		return ExpirePointResponse{}, fmt.Errorf(
			"decode expire-point response: %w",
			err,
		)
	}

	return gatewayResponse, nil
}

func (gateway PointGateway) transitionPoint(
	ctx context.Context,
	method string,
	path string,
	body PointTransitionRequest,
	action string,
) (Point, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return Point{}, fmt.Errorf(
			"marshal %s request: %w",
			action,
			err,
		)
	}

	endpoint := gateway.buildEndpoint(path)
	req, err := http.NewRequestWithContext(
		ctx,
		method,
		endpoint,
		bytes.NewReader(data),
	)
	if err != nil {
		return Point{}, fmt.Errorf(
			"create %s request: %w",
			action,
			err,
		)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	response, err := gateway.getHTTPClient().Do(req)
	if err != nil {
		return Point{}, fmt.Errorf(
			"call point-service %s endpoint: %w",
			action,
			err,
		)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(response.Body)
		return Point{}, fmt.Errorf(
			"point-service %s returned status %d: %s",
			action,
			response.StatusCode,
			string(responseBody),
		)
	}

	var gatewayResponse Point
	if err := json.NewDecoder(response.Body).Decode(
		&gatewayResponse,
	); err != nil {
		return Point{}, fmt.Errorf(
			"decode %s response: %w",
			action,
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
