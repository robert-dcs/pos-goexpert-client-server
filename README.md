# Go - Client & Server (Cotação do Dólar)

Projeto para praticar:

- Webserver HTTP
- Consumo de API externa
- Uso de `context` com timeout
- Persistência com SQLite
- Manipulação de arquivos

---

## Arquitetura

Client (executável) → Server (HTTP :8080) → API externa

Server:
GET http://localhost:8080/cotacao  
Retorno:
{ "bid": "5.12" }

Client:
- Executa uma vez
- Consome o Server
- Salva `cotacao.txt`
- Encerra execução

---

## Timeouts (Requisito)

| Componente | Timeout | Constante |
|------------|---------|-----------|
| Server → API | 200ms | `apiTimeout` |
| Server → DB | 10ms | `dbTimeout` |
| Client → Server | 300ms | `clientTimeout` |

Em caso de timeout:
- Erro é logado (`log.Println`)
- Server retorna erro HTTP
- Aplicação não é encerrada

---

## Execução

Instalar dependências:
go mod tidy

Subir o Server:
go run ./cmd/server

Executar o Client:
go run ./cmd/client

Arquivo gerado:
cotacao.txt

---

## Como Testar Cenários de Erro

Obs: linhas podem variar levemente dependendo da formatação.

| Cenário | Arquivo | Linha aprox. | Alteração |
|----------|----------|--------------|-----------|
| Timeout API | cmd/server/main.go | ~13 | `apiTimeout = 1 * time.Millisecond` |
| Timeout DB | cmd/server/main.go | ~14 | `dbTimeout = 1 * time.Nanosecond` |
| Timeout Client | cmd/client/main.go | ~12 | `clientTimeout = 1 * time.Millisecond` |
| Simular lentidão API | cmd/server/main.go | ~55 (dentro de `cotacaoHandler`) | Adicionar `time.Sleep(500 * time.Millisecond)` |
| API inválida | cmd/server/main.go | ~12 | Alterar `apiURL` para URL inexistente |
| Server indisponível | — | — | Executar client sem subir o server |

---

## Consultar Banco SQLite

Arquivo:
cotacoes.db

Comandos:

sqlite3 cotacoes.db  
.tables  
.mode column  
.headers on  
SELECT * FROM cotacoes;  
.exit  

---

## Stack

- Go 1.24
- net/http
- context
- SQLite (modernc.org/sqlite)
- encoding/json