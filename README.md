# Go - Client & Server (Cotação do Dólar)

Projeto em Go contendo dois executáveis (Server e Client) para praticar HTTP, consumo de API externa, uso de context com timeout, persistência em SQLite e manipulação de arquivos.

O Server consome a API https://economia.awesomeapi.com.br/json/last/USD-BRL, persiste o valor do dólar (USD → BRL) em um banco SQLite (`cotacoes.db`) e expõe o endpoint GET /cotacao na porta 8080 retornando apenas o campo `bid` em JSON. O Client consome esse endpoint e salva o valor recebido no arquivo `cotacao.txt` no formato: `Dólar: {valor}`.

Timeouts aplicados:
- Server → API externa: 200ms | Banco SQLite: 10ms
- Client → Requisição ao Server: 300ms
Todos utilizando `context` com propagação de cancelamento.

Estrutura:
cmd/server  → executável do servidor  
cmd/client  → executável do cliente  

Como executar:

1) Instalar dependências:
go mod tidy

2) Subir o Server:
go run ./cmd/server

Servidor disponível em:
http://localhost:8080/cotacao

3) Testar manualmente via terminal:
curl http://localhost:8080/cotacao

Ou testar via Postman:
- Método: GET
- URL: http://localhost:8080/cotacao

Resposta esperada:
{ "bid": "5.12" }

4) Rodar o Client (em outro terminal):
go run ./cmd/client

Arquivo gerado:
cotacao.txt

Stack utilizada:
- Go
- net/http
- context
- SQLite (github.com/mattn/go-sqlite3)
- encoding/json
- Manipulação de arquivos
