package handler

import (
	"errors"
	"testing"

	"sistema-notas/estoque/internal/core/domain"
)

func TestValidarNotaCalculaTotalNoBackend(t *testing.T) {
	nota := domain.NotaFiscal{ValorTotal: 999, Itens: []domain.ItemNota{{ProdutoCodigo: " PROD-1 ", Quantidade: 2, PrecoUnitario: 12.5}}}

	if err := validarNota(&nota); err != nil {
		t.Fatalf("validarNota() retornou erro: %v", err)
	}
	if nota.ValorTotal != 25 || nota.Itens[0].Subtotal != 25 || nota.Itens[0].ProdutoCodigo != "PROD-1" {
		t.Fatalf("nota normalizada incorretamente: %#v", nota)
	}
}

func TestValidarNotaRejeitaQuantidadeInvalida(t *testing.T) {
	err := validarNota(&domain.NotaFiscal{Itens: []domain.ItemNota{{ProdutoCodigo: "PROD-1", Quantidade: 0}}})
	if !errors.Is(err, domain.ErrValidacao) {
		t.Fatalf("erro = %v; esperado ErrValidacao", err)
	}
}
