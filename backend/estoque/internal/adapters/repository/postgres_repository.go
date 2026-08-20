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

// Salvar
func (r *ProdutoPostgres) Salvar(produto domain.Produto) (int, error) {
	sqlStatement := `INSERT INTO produtos (codigo, descricao, saldo) VALUES ($1, $2, $3) RETURNING id`

	var id int

	err := r.db.QueryRow(sqlStatement, produto.Codigo, produto.Descricao, produto.Saldo).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Listar
func (r *ProdutoPostgres) Listar() ([]domain.Produto, error) {
	query := "SELECT codigo, descricao, saldo FROM produtos ORDER BY codigo ASC"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var produtos []domain.Produto
	for rows.Next() {
		var produto domain.Produto
		if err := rows.Scan(&produto.Codigo, &produto.Descricao, &produto.Saldo); err != nil {
			return nil, err
		}
		produtos = append(produtos, produto)
	}

	return produtos, nil
}

func (r *ProdutoPostgres) BaixarEstoque(itens []domain.ItemNota) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := `UPDATE produtos SET saldo = saldo - $1 WHERE codigo = $2`

	for _, item := range itens {
		_, err := tx.Exec(query, item.Quantidade, item.ProdutoCodigo)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
