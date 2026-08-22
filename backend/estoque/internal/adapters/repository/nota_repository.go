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

func (r *NotaRepositoryPostgres) Emitir(nota domain.NotaFiscal) (domain.NotaFiscal, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return domain.NotaFiscal{}, err
	}

	queryNota := `INSERT INTO notas_fiscais (valor_total, status) VALUES ($1, 'ABERTA') RETURNING id, numero, data_emissao, status`
	err = tx.QueryRow(queryNota, nota.ValorTotal).Scan(&nota.ID, &nota.Numero, &nota.DataEmissao, &nota.Status)
	if err != nil {
		tx.Rollback()
		return domain.NotaFiscal{}, err
	}

	queryItem := `INSERT INTO itens_nota (nota_fiscal_id, produto_codigo, quantidade, preco_unitario, subtotal) VALUES ($1, $2, $3, $4, $5)`
	for _, item := range nota.Itens {
		_, err = tx.Exec(queryItem, nota.ID, item.ProdutoCodigo, item.Quantidade, item.PrecoUnitario, item.Subtotal)
		if err != nil {
			tx.Rollback()
			return domain.NotaFiscal{}, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return domain.NotaFiscal{}, err
	}

	return nota, nil
}

func (r *NotaRepositoryPostgres) Listar() ([]domain.NotaFiscal, error) {
	rows, err := r.db.Query(`SELECT id, numero, data_emissao, valor_total, status FROM notas_fiscais ORDER BY numero DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notas []domain.NotaFiscal
	for rows.Next() {
		var nota domain.NotaFiscal
		if err := rows.Scan(&nota.ID, &nota.Numero, &nota.DataEmissao, &nota.ValorTotal, &nota.Status); err != nil {
			return nil, err
		}
		notas = append(notas, nota)
	}
	return notas, rows.Err()
}

func (r *NotaRepositoryPostgres) BuscarNotaPorID(id int) (domain.NotaFiscal, error) {
	var nota domain.NotaFiscal
	query := `SELECT id, numero, data_emissao, valor_total, status FROM notas_fiscais WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&nota.ID, &nota.Numero, &nota.DataEmissao, &nota.ValorTotal, &nota.Status)
	if err != nil {
		return nota, err
	}

	queryItens := `SELECT produto_codigo, quantidade, preco_unitario, subtotal FROM itens_nota WHERE nota_fiscal_id = $1`
	rows, err := r.db.Query(queryItens, id)
	if err != nil {
		return nota, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.ItemNota
		if err := rows.Scan(&item.ProdutoCodigo, &item.Quantidade, &item.PrecoUnitario, &item.Subtotal); err != nil {
			return nota, err
		}
		nota.Itens = append(nota.Itens, item)
	}
	return nota, nil
}

func (r *NotaRepositoryPostgres) FecharNota(id int) error {
	result, err := r.db.Exec(`UPDATE notas_fiscais SET status = 'FECHADA' WHERE id = $1 AND status = 'ABERTA'`, id)
	if err != nil {
		return err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if changed == 0 {
		return domain.ErrNotaFechada
	}
	return nil
}
