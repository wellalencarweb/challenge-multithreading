package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
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

func printTable(api string, data map[string]interface{}) {
	fmt.Printf("Resposta da API: %s\n\n", api)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for k, v := range data {
		// Apenas campos úteis
		if strings.Contains(k, "cep") || strings.Contains(k, "logradouro") || strings.Contains(k, "bairro") || strings.Contains(k, "localidade") || strings.Contains(k, "uf") || strings.Contains(k, "city") || strings.Contains(k, "state") {
			fmt.Fprintf(w, "%s:\t%v\n", strings.Title(k), v)
		}
	}
	w.Flush()
}

func isValidCep(cep string) bool {
	if len(cep) != 8 {
		return false
	}
	for _, c := range cep {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <CEP>")
		return
	}

	cep := os.Args[1]
	cep = strings.ReplaceAll(cep, "-", "") // remover hífen

	if !isValidCep(cep) {
		fmt.Println("CEP inválido. Informe apenas números (8 dígitos).")
		return
	}

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
		printTable(res.API, res.Data)
	case <-ctx.Done():
		fmt.Println("⏰ Timeout: nenhuma API respondeu dentro de 1 segundo.")
	}
}
