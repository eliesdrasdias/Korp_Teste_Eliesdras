# Sistema de Emissão de Notas Fiscais

Teste técnico com Angular 18, Go e PostgreSQL. Há dois processos HTTP independentes: **Estoque** (`:8080`), para produtos e saldos, e **Faturamento** (`:8081`), para notas e seu fechamento.

## Executando

1. Copie `.env.example` para `.env` e ajuste a senha local.
2. Inicie o banco: `docker compose up -d postgres`.
3. Exporte `DATABASE_URL` conforme `.env` e execute `go run ./cmd/migration` em `backend/estoque`.
4. Em terminais separados, execute `go run ./cmd/estoque` e `go run ./cmd/faturamento`.
5. Em `frontend`, execute `npm install` e `npm start`; acesse `http://localhost:4200`.

`DATABASE_URL` e `ESTOQUE_URL` são variáveis de ambiente. `.env` não deve ser commitido.

## Contratos HTTP

| Serviço | Endpoint | Uso |
| --- | --- | --- |
| Estoque | `POST /produtos` | cria produto; código único e saldo não negativo |
| Estoque | `GET /produtos/listar` | lista produtos e saldo atual |
| Estoque | `POST /produtos/baixa` | endpoint interno usado pelo Faturamento |
| Faturamento | `POST /notas` | cria nota `ABERTA`; número vem do banco |
| Faturamento | `GET /notas/listar` | lista notas |
| Faturamento | `POST /notas/imprimir` | fecha uma nota aberta e baixa o estoque |

Os erros possuem JSON consistente com `status`, `message`, `timestamp` e `details`. Regras de negócio retornam `422`, repetição do fechamento retorna `409` e indisponibilidade do Estoque retorna `503`.

## Decisões para explicar na entrevista

- A sequência da nota é PostgreSQL (`numero SERIAL`), não do frontend.
- A nota inicia `ABERTA`; o update `WHERE status = 'ABERTA'` impede baixa duplicada por repetição de fechamento.
- A baixa consolida itens, usa transação e `SELECT ... FOR UPDATE`, valida o saldo e só então atualiza: notas concorrentes não consomem o mesmo saldo.
- O Faturamento usa timeout de quatro segundos ao chamar o Estoque. Falha de rede retorna mensagem amigável e mantém a nota aberta.
- Não existe transação distribuída: inserir saga/outbox seria desproporcional para o teste. O risco residual de falha após a baixa do Estoque é documentado; em produção seria resolvido por outbox/reconciliação. A operação no Estoque é atômica.
- Go usa `go.mod`/`go.sum`. Angular usa standalone components, `HttpClient`, Reactive Forms para produto, `ngOnInit` para carga inicial, `takeUntilDestroyed` para subscriptions e `finalize` para loading.

## Roteiro manual

1. Cadastre `PROD001`, `Notebook`, saldo `10`.
2. Crie nota com `PROD001`, quantidade `2` e preço `100`.
3. Clique em **Imprimir Nota**: status `FECHADA`, saldo `8`.
4. Tente imprimir novamente: o botão bloqueia e uma chamada direta recebe `409`.
5. Crie nota de quantidade `9`: recebe `422`; a nota fica aberta.
6. Pare o Estoque e feche uma nota: a interface informa indisponibilidade; nada é fechado.

## Verificação

`npm run build` foi executado com sucesso. Para validar Go localmente, execute `gofmt -w cmd internal` e `go test ./...` em `backend/estoque`; neste ambiente o binário Go instalado via Snap é bloqueado pelo sandbox.
