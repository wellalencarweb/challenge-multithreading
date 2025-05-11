package fetchers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wellalencarweb/challenge-multithreading/domain"
)

func FindByViaCEP(ctx context.Context, cep string, ch chan<- domain.Result) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- domain.Result{API: "ViaCEP", Error: err}
		return
	}
	defer resp.Body.Close()

	var raw struct {
		Cep          string `json:"cep"`
		State        string `json:"uf"`
		City         string `json:"localidade"`
		Neighborhood string `json:"bairro"`
		Street       string `json:"logradouro"`
	}

	err = json.NewDecoder(resp.Body).Decode(&raw)
	endereco := domain.Address{
		Cep:          raw.Cep,
		State:        raw.State,
		City:         raw.City,
		Neighborhood: raw.Neighborhood,
		Street:       raw.Street,
	}

	ch <- domain.Result{API: "ViaCEP", Data: endereco, Error: err}
}
