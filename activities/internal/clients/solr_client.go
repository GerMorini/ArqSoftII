package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	dto "activities/internal/dto"
)

// IndexActivityToSolr indexa solamente los campos permitidos por el usuario
func IndexActivityToSolr(solrURL, collection string, a dto.ActivityAdministration) error {
	doc := map[string]interface{}{
		"id":          a.ID,
		"nombre":      a.Nombre,
		"descripcion": a.Descripcion,
		"profesor":    a.Profesor,
		"dia_semana":  a.DiaSemana,
	}

	payload := map[string]interface{}{"add": []interface{}{doc}}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("%s/%s/update?commit=true", solrURL, collection)
	resp, err := client.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("solr returned status %d", resp.StatusCode)
	}
	return nil
}

func DeleteActivityFromSolr(solrURL, collection, id string) error {
	payload := map[string]interface{}{"delete": map[string]string{"id": id}}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("%s/%s/update?commit=true", solrURL, collection)
	resp, err := client.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("solr returned status %d", resp.StatusCode)
	}
	return nil
}
