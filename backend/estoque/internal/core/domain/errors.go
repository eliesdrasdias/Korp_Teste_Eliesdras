package domain

import "errors"

var (
	ErrValidacao          = errors.New("dados inválidos")
	ErrProdutoInexistente = errors.New("produto não encontrado")
	ErrSaldoInsuficiente  = errors.New("saldo insuficiente")
	ErrNotaFechada        = errors.New("nota já está fechada")
	ErrNotaSemItens       = errors.New("a nota precisa ter pelo menos um item")
)
