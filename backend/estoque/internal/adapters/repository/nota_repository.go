package repository

import (
	"database/sql"
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

	queryItem := `INSERT INTO itens_nota (nota_fiscal_id, produto_codigo, quantidade, preco_unitario, subtotal) VALUES ($1, $2, $3, $4, $5)`
	for _, item := range nota.Itens {
		_, err = tx.Exec(queryItem, notaID, item.ProdutoCodigo, item.Quantidade, item.PrecoUnitario, item.Subtotal)
		if err != nil {
			tx.Rollback()
			return 0, err
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

func (r *NotaRepositoryPostgres) BuscarNotaPorID(id int) (domain.NotaFiscal, error) {
	var nota domain.NotaFiscal
	query := `SELECT id, status FROM notas_fiscais WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&nota.ID, &nota.Status)
	if err != nil {
		return nota, err
	}

	queryItens := `SELECT produto_codigo, quantidade FROM itens_nota WHERE nota_fiscal_id = $1`
	rows, err := r.db.Query(queryItens, id)
	if err != nil {
		return nota, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.ItemNota
		rows.Scan(&item.ProdutoCodigo, &item.Quantidade)
		nota.Itens = append(nota.Itens, item)
	}
	return nota, nil
}

func (r *NotaRepositoryPostgres) FecharNota(id int) error {
	_, err := r.db.Exec(`UPDATE notas_fiscais SET status = 'Fechada' WHERE id = $1`, id)
	return err
}
