package ports

import "sistema-notas/estoque/internal/core/domain"

type ProdutoRepository interface {
	Salvar(produto domain.Produto) (int, error)
	Listar() ([]domain.Produto, error)
	BaixarEstoque(itens []domain.ItemNota) error
}

type NotaRepository interface {
	Emitir(nota domain.NotaFiscal) (domain.NotaFiscal, error)
	Listar() ([]domain.NotaFiscal, error)
	BuscarNotaPorID(id int) (domain.NotaFiscal, error)
	FecharNota(id int) error
}
