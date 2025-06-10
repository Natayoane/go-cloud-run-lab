package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
)

type ViaCepData struct {
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
}

func ValidateCEP(cep string) error {
	cleaned := strings.ReplaceAll(cep, "-", "")
	if len(cleaned) != 8 || !regexp.MustCompile(`^[0-9]+$`).MatchString(cleaned) {
		return fmt.Errorf("invalid zipcode: %s", cep)
	}
	return nil
}

func GetViaCepAPI(ctx context.Context, cep string) (ViaCepData, error) {
	var result ViaCepData
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	log.Printf("CEP API URL: %s", url)
	
	body, err := FetchAPI(ctx, url)
	if err != nil {
		log.Printf("Error fetching CEP data: %v", err)
		return ViaCepData{}, err
	}
	
	log.Printf("CEP API Response: %s", body.body)
	
	if err := json.Unmarshal([]byte(body.body), &result); err != nil {
		log.Printf("Error unmarshaling CEP data: %v", err)
		return ViaCepData{}, err
	}
	
	log.Printf("CEP Data: %+v", result)
	return result, nil
} 