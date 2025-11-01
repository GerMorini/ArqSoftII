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

	log "github.com/sirupsen/logrus"
)

type SolrClient struct {
	baseURL string
	core    string
	client  *http.Client
}

type SolrDocument struct {
	ID          string   `json:"id"`
	Titulo      []string `json:"titulo"`
	Descripcion []string `json:"descripcion"`
	DiaSemana   []string `json:"dia"`
}

type SolrResponse struct {
	Response struct {
		NumFound int            `json:"numFound"`
		Start    int            `json:"start"`
		Docs     []SolrDocument `json:"docs"`
	} `json:"response"`
}

type SolrUpdateResponse struct {
	ResponseHeader struct {
		Status int `json:"status"`
		QTime  int `json:"QTime"`
	} `json:"responseHeader"`
}

const (
	defaultCount = 10
)

func NewSolrClient(host, port, core string) *SolrClient {
	baseURL := fmt.Sprintf("http://%s:%s/solr/%s", host, port, core)
	return &SolrClient{
		baseURL: baseURL,
		core:    core,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *SolrClient) Index(ctx context.Context, activity dto.Activity) error {
	doc := SolrDocument{
		ID:          activity.ID,
		Titulo:      []string{activity.Titulo},
		Descripcion: []string{activity.Descripcion},
		DiaSemana:   []string{activity.DiaSemana},
	}

	data, err := json.Marshal([]SolrDocument{doc})
	if err != nil {
		return fmt.Errorf("error marshalling document: %w", err)
	}

	url := fmt.Sprintf("%s/update?commit=true", s.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(data)))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("solr returned status %d", resp.StatusCode)
	}

	var updateResp SolrUpdateResponse
	if err := json.NewDecoder(resp.Body).Decode(&updateResp); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	if updateResp.ResponseHeader.Status != 0 {
		return fmt.Errorf("solr update failed with status %d", updateResp.ResponseHeader.Status)
	}

	return nil
}

func (s *SolrClient) Search(ctx context.Context, query string, page int, count int) (dto.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if count <= 0 {
		count = defaultCount
	}

	// calcular offset
	start := (page - 1) * count

	params := url.Values{}
	params.Set("q", query)
	params.Set("wt", "json")
	params.Set("start", fmt.Sprintf("%d", start))
	params.Set("rows", fmt.Sprintf("%d", count))

	url := fmt.Sprintf("%s/select?%s", s.baseURL, params.Encode())
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return dto.PaginatedResponse{}, fmt.Errorf("error creating request: %w", err)
	}
	log.Infof("REQUEST: %v", req)

	resp, err := s.client.Do(req)
	if err != nil {
		return dto.PaginatedResponse{}, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dto.PaginatedResponse{}, fmt.Errorf("solr returned status %d", resp.StatusCode)
	}

	var solrResp SolrResponse
	if err := json.NewDecoder(resp.Body).Decode(&solrResp); err != nil {
		return dto.PaginatedResponse{}, fmt.Errorf("error decoding response: %w", err)
	}

	activities := make([]dto.Activity, len(solrResp.Response.Docs))
	for i, doc := range solrResp.Response.Docs {
		activities[i] = dto.Activity{
			ID:          doc.ID,
			Titulo:      doc.Titulo[0],
			Descripcion: doc.Descripcion[0],
			DiaSemana:   doc.DiaSemana[0],
		}
	}

	return dto.PaginatedResponse{
		Page:    page,
		Count:   len(activities),
		Total:   solrResp.Response.NumFound,
		Results: activities,
	}, nil
}

func (s *SolrClient) Delete(ctx context.Context, id string) error {
	data := fmt.Sprintf(`{"delete":{"id":"%s"}}`, id)
	url := fmt.Sprintf("%s/update?commit=true", s.baseURL)

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(data))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("solr returned status %d", resp.StatusCode)
	}

	var updateResp SolrUpdateResponse
	if err := json.NewDecoder(resp.Body).Decode(&updateResp); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	if updateResp.ResponseHeader.Status != 0 {
		return fmt.Errorf("solr delete failed with status %d", updateResp.ResponseHeader.Status)
	}

	return nil
}

func (s *SolrClient) Commit(ctx context.Context) error {
	url := fmt.Sprintf("%s/update?commit=true", s.baseURL)

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(""))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("solr returned status %d", resp.StatusCode)
	}

	var updateResp SolrUpdateResponse
	if err := json.NewDecoder(resp.Body).Decode(&updateResp); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	if updateResp.ResponseHeader.Status != 0 {
		return fmt.Errorf("solr commit failed with status %d", updateResp.ResponseHeader.Status)
	}

	return nil
}
