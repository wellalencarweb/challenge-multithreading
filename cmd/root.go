package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/wellalencarweb/challenge-multithreading/service"
)

func Run() {
	app := &cli.App{
		Name:  "cepfinder",
		Usage: "Busca de CEPs via linha de comando com concorrência",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "cep", Aliases: []string{"c"}, Required: true},
		},
		Action: func(c *cli.Context) error {
			cep := c.String("cep")

			cep = strings.ReplaceAll(cep, "-", "")

			if !isValidCep(cep) {
				return fmt.Errorf("❌ CEP inválido. Informe apenas números (8 dígitos)")
			}

			res, err := service.FindCep(c.Context, cep)
			if err != nil {
				return fmt.Errorf("Erro: %w", err)
			}

			fmt.Printf("\n✅ API: %s\n", res.API)
			fmt.Printf("📦 Cep: %s\n", res.Data.Cep)
			fmt.Printf("🗺️  State: %s\n", res.Data.State)
			fmt.Printf("🏙️  City: %s\n", res.Data.City)
			fmt.Printf("🏘️  Neighborhood: %s\n", res.Data.Neighborhood)
			fmt.Printf("🛣️  Street: %s\n", res.Data.Street)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func isValidCep(cep string) bool {
	re := regexp.MustCompile(`^\d{8}$`)
	return re.MatchString(cep)
}
