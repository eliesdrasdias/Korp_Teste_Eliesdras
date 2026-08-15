package repository

import (
	"database/sql"
	"sistema-notas/estoque/internal/core/domain"
)

type ProdutoPostgres struct {
	db *sql.DB
}

func NewProdutoPostgres(db *sql.DB) *ProdutoPostgres {
	return &ProdutoPostgres{db: db}
}

func (r *ProdutoPostgres) Salvar(produto domain.Produto) (int, error) {
	sqlStatement := `INSERT INTO produtos (codigo, descricao, saldo) VALUES ($1, $2, $3) RETURNING id`

	var id int

	err := r.db.QueryRow(sqlStatement, produto.Codigo, produto.Descricao, produto.Saldo).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
