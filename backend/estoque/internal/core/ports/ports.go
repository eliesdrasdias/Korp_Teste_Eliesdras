package ports

import "sistema-notas/estoque/internal/core/domain"

type ProdutoRepository interface {
	Salvar(produto domain.Produto) (int, error)
}
