# 🧾 Sistema de Emissão de Notas Fiscais e Controle de Estoque

[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Angular](https://img.shields.io/badge/Angular-18-DD0031?logo=angular&logoColor=white)](https://angular.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white)](https://www.docker.com/)

Aplicação de demonstração para emissão de notas fiscais e controle de estoque. O projeto foi construído para evidenciar decisões práticas de arquitetura: separação clara de responsabilidades, persistência relacional, validação de regras críticas no backend e comunicação resiliente entre serviços.

O usuário cadastra produtos, monta uma nota com múltiplos itens e a fecha por meio da ação **Imprimir Nota**. Nesse momento, o sistema valida o estoque novamente, realiza a baixa e altera o status da nota de `ABERTA` para `FECHADA`.

## ✨ Tecnologias utilizadas

| Camada | Tecnologias |
| --- | --- |
| Frontend | Angular 18, TypeScript, Reactive Forms, HttpClient e RxJS |
| Backend | Go, `net/http`, `database/sql`, `lib/pq` e Goose migrations |
| Banco de dados | PostgreSQL 15 |
| Infraestrutura local | Docker e Docker Compose |

## 🏗️ Arquitetura do projeto

O backend aplica uma organização inspirada em **Arquitetura Hexagonal (Ports & Adapters)**. O domínio contém entidades e contratos; os adaptadores implementam a comunicação HTTP e a persistência PostgreSQL.

```text
Angular (porta 4200)
       │
       ├──────────────► Estoque (porta 8080)
       │                    └── Produtos, saldos e baixa de estoque
       │
       └──────────────► Faturamento (porta 8081)
                            └── Notas fiscais e fechamento
                                     │
                                     └── HTTP ──► Estoque

                         PostgreSQL (porta 5432)
```

### Microsserviço de Estoque — `:8080`

Responsável por cadastrar e listar produtos, validar disponibilidade e baixar saldo. A baixa é executada em transação: os produtos são bloqueados no banco, o saldo é conferido e somente então atualizado. Isso evita saldo negativo e reduz riscos de concorrência.

### Microsserviço de Faturamento — `:8081`

Responsável por criar, listar e fechar notas fiscais. A numeração é gerada pelo PostgreSQL e o status inicial é `ABERTA`. No fechamento, o serviço solicita a baixa ao Estoque e somente após a confirmação altera a nota para `FECHADA`.

### Ports & Adapters no Go

```text
internal/
├── core/
│   ├── domain/       # Entidades e erros de negócio
│   └── ports/        # Interfaces de repositório
└── adapters/
    ├── handler/      # Adaptadores HTTP / controllers
    └── repository/   # Adaptadores PostgreSQL
```

Essa divisão permite que regras de negócio e contratos sejam menos acoplados ao protocolo HTTP e ao banco de dados.

## 🛡️ Destaque técnico: resiliência entre serviços

O endpoint de impressão no Faturamento chama o Estoque por HTTP com timeout de quatro segundos. Se o Estoque estiver indisponível, não responder a tempo ou retornar um erro inesperado:

- o detalhe técnico é registrado nos logs;
- o Faturamento responde `503 Service Unavailable` com mensagem amigável;
- a aplicação Angular apresenta o erro sem quebrar a tela;
- a nota permanece `ABERTA`, pois o fechamento só ocorre depois da confirmação da baixa.

Exemplo de mensagem exibida:

> Não foi possível concluir a nota porque o serviço de estoque está temporariamente indisponível. Tente novamente em alguns instantes.

## 📁 Estrutura de diretórios

```text
.
├── backend/
│   └── estoque/
│       ├── cmd/
│       │   ├── estoque/          # Inicialização do microsserviço de Estoque (:8080)
│       │   ├── faturamento/      # Inicialização do microsserviço de Faturamento (:8081)
│       │   └── migration/        # Executor das migrations Goose
│       ├── db/migrations/        # Schema PostgreSQL versionado
│       ├── internal/
│       │   ├── adapters/
│       │   │   ├── handler/      # Rotas HTTP, validações e respostas de erro
│       │   │   └── repository/   # Queries e transações PostgreSQL
│       │   ├── config/           # Leitura de variáveis de ambiente
│       │   └── core/
│       │       ├── domain/       # Produto, nota, item e erros de domínio
│       │       └── ports/        # Interfaces de persistência
│       ├── go.mod
│       └── go.sum
├── frontend/
│   ├── src/app/
│   │   ├── components/           # Cabeçalho e cadastro de produtos
│   │   ├── notas/                # Tela de emissão, itens e impressão de nota
│   │   └── services/             # Clientes HTTP dos dois serviços
│   ├── angular.json
│   └── package.json
├── docker-compose.yml            # PostgreSQL local
└── README.md
```

## ✅ Pré-requisitos

Antes de iniciar, instale:

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) ou Docker Engine com Docker Compose v2;
- [Go](https://go.dev/dl/) 1.25 ou compatível com o `go.mod`;
- [Node.js](https://nodejs.org/) 20 LTS ou superior;
- npm (incluído com Node.js);
- Angular CLI é opcional, pois o projeto executa o CLI local via `npm start`.

## 🚀 Como executar localmente

### 1. Clone e acesse o projeto

```bash
git clone <URL_DO_SEU_REPOSITORIO>
cd Korp_Test_Eliesdras
```

### 2. Suba o PostgreSQL com Docker

Na raiz do projeto:

```bash
docker compose up -d postgres
docker compose ps
```

O PostgreSQL ficará disponível em `localhost:5432`, com banco `sistema_notas`. As credenciais padrão do Compose são `postgres` / `postgres`.

### 3. Execute as migrations

Abra um terminal na pasta do backend e exporte a connection string:

```bash
cd backend/estoque
export DATABASE_URL='host=localhost port=5432 user=postgres password=postgres dbname=sistema_notas sslmode=disable'
go run ./cmd/migration
```

No Windows PowerShell, use:

```powershell
$env:DATABASE_URL = 'host=localhost port=5432 user=postgres password=postgres dbname=sistema_notas sslmode=disable'
go run ./cmd/migration
```

### 4. Inicie o microsserviço de Estoque

Em um terminal, dentro de `backend/estoque`:

```bash
export DATABASE_URL='host=localhost port=5432 user=postgres password=postgres dbname=sistema_notas sslmode=disable'
go run ./cmd/estoque
```

O serviço estará em `http://localhost:8080`.

### 5. Inicie o microsserviço de Faturamento

Em outro terminal, dentro de `backend/estoque`:

```bash
export DATABASE_URL='host=localhost port=5432 user=postgres password=postgres dbname=sistema_notas sslmode=disable'
export ESTOQUE_URL='http://localhost:8080'
go run ./cmd/faturamento
```

O serviço estará em `http://localhost:8081`.

### 6. Inicie o frontend Angular

Em um terceiro terminal, na raiz do projeto:

```bash
cd frontend
npm install
npm start
```

Abra [http://localhost:4200](http://localhost:4200) no navegador.

## 🔎 Fluxo rápido de demonstração

1. Cadastre um produto: código `PROD001`, descrição `Notebook`, saldo `10`.
2. Na tela de notas, selecione o produto, informe quantidade `2` e preço unitário.
3. Clique em **Adicionar** e depois em **Criar nota fiscal**.
4. Clique em **Imprimir Nota** para fechar a nota.
5. Retorne aos produtos e confirme que o saldo passou de `10` para `8`.

Para validar a resiliência, interrompa o processo de Estoque e tente imprimir uma nova nota aberta. O Faturamento retornará `503` e a nota permanecerá aberta.

## 🧪 Comandos úteis

```bash
# Build do frontend
cd frontend && npm run build

# Testes do backend
cd backend/estoque && go test ./...

# Consultas rápidas às APIs
curl http://localhost:8080/produtos/listar
curl http://localhost:8081/notas/listar
```

---

Desenvolvido como teste técnico com foco em clareza arquitetural, regras de negócio no backend e uma experiência de demonstração objetiva.
