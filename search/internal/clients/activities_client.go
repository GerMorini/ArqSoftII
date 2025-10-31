package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"search/internal/dto"
	"strings"
	"time"
)

type ActivitiesClient struct {
	baseURL string
	client  *http.Client
}

type ActivitiesAPIResponse struct {
	Activities []dto.Activity `json:"activities"`
	Count      int            `json:"count"`
}

type ActivityAPIResponse struct {
	Activity dto.Activity `json:"activity"`
}

func NewActivitiesClient(baseURL string) *ActivitiesClient {
	return &ActivitiesClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *ActivitiesClient) GetActivityByID(ctx context.Context, id string) (*dto.Activity, error) {
	url := fmt.Sprintf("%s/activities/%s", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("activity not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("activities API returned status %d", resp.StatusCode)
	}

	var apiResp ActivityAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &apiResp.Activity, nil
}

func (c *ActivitiesClient) GetActivitiesByIDs(ctx context.Context, ids []string) ([]dto.Activity, error) {
	if len(ids) == 0 {
		return []dto.Activity{}, nil
	}

	// Build query parameter: /activities/many?ids=id1,id2,id3
	idsParam := url.QueryEscape(strings.Join(ids, ","))
	requestURL := fmt.Sprintf("%s/activities/many?ids=%s", c.baseURL, idsParam)

	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("activities API returned status %d", resp.StatusCode)
	}

	var apiResp ActivitiesAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return apiResp.Activities, nil
}

func (c *ActivitiesClient) GetAllActivities(ctx context.Context) ([]dto.Activity, error) {
	requestURL := fmt.Sprintf("%s/activities", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("activities API returned status %d", resp.StatusCode)
	}

	var activities []dto.Activity
	if err := json.NewDecoder(resp.Body).Decode(&activities); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return activities, nil
}
