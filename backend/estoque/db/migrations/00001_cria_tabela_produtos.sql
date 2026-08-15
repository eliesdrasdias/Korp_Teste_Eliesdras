-- +goose Up
CREATE TABLE IF NOT EXISTS produtos (
    id SERIAL PRIMARY KEY,
    codigo VARCHAR(50) UNIQUE NOT NULL,
    descricao TEXT NOT NULL,
    saldo NUMERIC(10, 2) NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS produtos;