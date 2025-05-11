package fetchers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wellalencarweb/challenge-multithreading/domain"
)

func FindByBrasilAPI(ctx context.Context, cep string, ch chan<- domain.Result) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- domain.Result{API: "BrasilAPI", Error: err}
		return
	}
	defer resp.Body.Close()

	var endereco domain.Address
	err = json.NewDecoder(resp.Body).Decode(&endereco)
	ch <- domain.Result{API: "BrasilAPI", Data: endereco, Error: err}
}
