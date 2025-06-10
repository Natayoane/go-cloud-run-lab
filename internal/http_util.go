package internal

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

type apiResponse struct {
	body       string
	statusCode int
	err        error
}

func FetchAPI(ctx context.Context, url string) (apiResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return apiResponse{}, err
	}
	
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return apiResponse{}, err
	}
	defer res.Body.Close()
	
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return apiResponse{}, err
	}
	
	log.Printf("API Response Status: %d", res.StatusCode)
	log.Printf("API Response Body: %s", string(body))
	
	if res.StatusCode != http.StatusOK {
		return apiResponse{}, fmt.Errorf("API returned status %d: %s", res.StatusCode, string(body))
	}
	
	return apiResponse{
		body:       string(body),
		statusCode: res.StatusCode,
	}, nil
} 