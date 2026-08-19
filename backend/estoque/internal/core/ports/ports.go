package ports

import "sistema-notas/estoque/internal/core/domain"

type ProdutoRepository interface {
	Salvar(produto domain.Produto) (int, error)
	Listar() ([]domain.Produto, error)
}

type NotaRepository interface {
	Emitir(nota domain.NotaFiscal) (int, error)
	Listar() ([]domain.NotaFiscal, error)
}
