package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Result struct {
	API   string
	Data  map[string]interface{}
	Error error
}

func fetchFromBrasilAPI(ctx context.Context, cep string, ch chan<- Result) {
	url := "https://brasilapi.com.br/api/cep/v1/" + cep
	resp, err := httpRequestWithContext(ctx, url)
	if err != nil {
		ch <- Result{API: "BrasilAPI", Error: err}
		return
	}
	ch <- Result{API: "BrasilAPI", Data: resp}
}

func fetchFromViaCEP(ctx context.Context, cep string, ch chan<- Result) {
	url := "http://viacep.com.br/ws/" + cep + "/json/"
	resp, err := httpRequestWithContext(ctx, url)
	if err != nil {
		ch <- Result{API: "ViaCEP", Error: err}
		return
	}
	ch <- Result{API: "ViaCEP", Data: resp}
}

func httpRequestWithContext(ctx context.Context, url string) (map[string]interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func main() {
	cep := "01425080"
	ch := make(chan Result, 2)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	go fetchFromBrasilAPI(ctx, cep, ch)
	go fetchFromViaCEP(ctx, cep, ch)

	select {
	case res := <-ch:
		if res.Error != nil {
			fmt.Printf("Erro na requisição da %s: %v\n", res.API, res.Error)
			return
		}
		fmt.Printf("Resposta da API: %s\n", res.API)
		for k, v := range res.Data {
			fmt.Printf("%s: %v\n", k, v)
		}
	case <-ctx.Done():
		fmt.Println("Timeout: nenhuma API respondeu dentro de 1 segundo.")
	}
}
