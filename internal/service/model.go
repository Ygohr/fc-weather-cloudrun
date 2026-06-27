package service

import "encoding/json"

type viaCEPResponse struct {
	Localidade string          `json:"localidade"`
	Erro       json.RawMessage `json:"erro"`
}

func (r viaCEPResponse) IsNotFound() bool {
	if len(r.Erro) == 0 {
		return false
	}

	var boolValue bool
	if err := json.Unmarshal(r.Erro, &boolValue); err == nil {
		return boolValue
	}

	var stringValue string
	if err := json.Unmarshal(r.Erro, &stringValue); err == nil {
		return stringValue == "true"
	}

	return false
}

type weatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}
