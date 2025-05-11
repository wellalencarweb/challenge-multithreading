# 🧵 challenge-multithreading — Buscador concorrente de CEPs

Este projeto é um **CLI (Command Line Interface)** em Go que realiza consultas de CEPs simultaneamente em duas APIs públicas (BrasilAPI e ViaCEP), retornando a resposta mais rápida entre elas. Utiliza concorrência (`goroutines`, `context`, `select`) com timeout de 1 segundo.

## 🚀 Funcionalidades

- Consulta simultânea em:
  - [`https://brasilapi.com.br`](https://brasilapi.com.br)
  - [`http://viacep.com.br`](http://viacep.com.br)
- Exibe apenas a resposta mais rápida.
- Timeout configurado para 1 segundo.
- Validação de CEP.
- Código organizado em estrutura Clean Architecture.

## 📦 Estrutura do Projeto

```
cepfinder/
├── cmd/              # CLI principal
├── fetchers/         # Comunicação com APIs externas
├── domain/           # Structs e modelos de dados
├── service/          # Lógica para decidir a resposta mais rápida
├── utils/            # Funções auxiliares (timeout, etc)
├── main.go           # Inicializador da aplicação
└── go.mod
```

## 🛠️ Requisitos

- Go 1.18 ou superior instalado.

## 🧪 Instalação e uso

### 1. Clone o repositório

```bash
git clone https://github.com/wellalencarweb/challenge-multithreading.git
cd challenge-multithreading
```

### 2. Compile o binário

```bash
go build -o cepfinder
```

### 3. Execute com um CEP válido

```bash
./cepfinder --cep 01153000
```

Ou com atalho:

```bash
./cepfinder -c 01153000
```

## 💡 Exemplo de saída

```
✅ API: BrasilAPI
📦 Cep: 01153000
🗺️  State: SP
🏙️  City: São Paulo
🏘️  Neighborhood: Barra Funda
🛣️  Street: Rua Vitorino Carmilo
```
▶️ Execução direta via main.go (sem binário)
Você também pode executar diretamente o projeto com o comando:
```bash
go run main.go --cep 01153000
```

Ou com atalho:
```bash
go run main.go --c 01153000
```

## 🧠 Conceitos aplicados

- `goroutines` para chamadas simultâneas.
- `context.WithTimeout()` para controlar o tempo limite.
- Canal (`chan`) para coletar o primeiro resultado com `select`.
- `urfave/cli` para parsing e flags de linha de comando.
- Clean code e separação de responsabilidades.