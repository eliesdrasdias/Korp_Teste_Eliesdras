-- +goose Up
ALTER TABLE produtos ADD CONSTRAINT produtos_saldo_nao_negativo CHECK (saldo >= 0);

CREATE TABLE IF NOT EXISTS notas_fiscais (
    id SERIAL PRIMARY KEY,
    numero SERIAL UNIQUE NOT NULL,
    data_emissao TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    valor_total NUMERIC(12, 2) NOT NULL DEFAULT 0,
    status VARCHAR(10) NOT NULL DEFAULT 'ABERTA' CHECK (status IN ('ABERTA', 'FECHADA'))
);

CREATE TABLE IF NOT EXISTS itens_nota (
    id SERIAL PRIMARY KEY,
    nota_fiscal_id INTEGER NOT NULL REFERENCES notas_fiscais(id),
    produto_codigo VARCHAR(50) NOT NULL,
    quantidade INTEGER NOT NULL CHECK (quantidade > 0),
    preco_unitario NUMERIC(12, 2) NOT NULL CHECK (preco_unitario >= 0),
    subtotal NUMERIC(12, 2) NOT NULL CHECK (subtotal >= 0)
);

-- +goose Down
DROP TABLE IF EXISTS itens_nota;
DROP TABLE IF EXISTS notas_fiscais;
ALTER TABLE produtos DROP CONSTRAINT IF EXISTS produtos_saldo_nao_negativo;
