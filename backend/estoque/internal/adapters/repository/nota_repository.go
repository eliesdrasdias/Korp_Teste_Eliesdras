package repository

import (
	"database/sql"
	"errors"
	"sistema-notas/estoque/internal/core/domain"
)

type NotaRepositoryPostgres struct {
	db *sql.DB
}

func NewNotaRepositoryPostgres(db *sql.DB) *NotaRepositoryPostgres {
	return &NotaRepositoryPostgres{db: db}
}

func (r *NotaRepositoryPostgres) Emitir(nota domain.NotaFiscal) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	queryNota := `INSERT INTO notas_fiscais (valor_total) VALUES ($1) RETURNING id`
	var notaID int
	err = tx.QueryRow(queryNota, nota.ValorTotal).Scan(&notaID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	queryItem := `INSERT INTO itens_nota (nota_id, produto_id, quantidade, preco_unitario, subtotal) VALUES ($1, $2, $3, $4, $5)`
	for _, item := range nota.Itens {
		_, err = tx.Exec(queryItem, notaID, item.ProdutoCodigo, item.Quantidade, item.PrecoUnitario, item.Subtotal)
		if err != nil {
			tx.Rollback()
			return 0, errors.New("Erro ao salvar os itens da nota")
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return notaID, nil
}

func (r *NotaRepositoryPostgres) Listar() ([]domain.NotaFiscal, error) {
	return nil, nil
}
