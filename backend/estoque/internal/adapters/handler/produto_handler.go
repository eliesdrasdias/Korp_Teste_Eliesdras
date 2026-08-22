package handler

import (
	"errors"
	"net/http"
	"strings"

	"sistema-notas/estoque/internal/core/domain"
	"sistema-notas/estoque/internal/core/ports"
)

type ProdutoHandler struct{ repo ports.ProdutoRepository }

func NewProdutoHandler(repo ports.ProdutoRepository) *ProdutoHandler {
	return &ProdutoHandler{repo: repo}
}

func (h *ProdutoHandler) Criar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}
	var produto domain.Produto
	if err := decodeJSON(r, &produto); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	produto.Codigo, produto.Descricao = strings.TrimSpace(produto.Codigo), strings.TrimSpace(produto.Descricao)
	if produto.Codigo == "" || produto.Descricao == "" || produto.Saldo < 0 {
		writeError(w, http.StatusUnprocessableEntity, "Código e descrição são obrigatórios e o saldo não pode ser negativo")
		return
	}
	id, err := h.repo.Salvar(produto)
	if err != nil {
		logUnexpected("falha ao salvar produto", err)
		if strings.Contains(err.Error(), "duplicate key") {
			writeError(w, http.StatusConflict, "Já existe um produto com este código")
			return
		}
		writeError(w, http.StatusInternalServerError, "Não foi possível cadastrar o produto")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"message": "Produto cadastrado com sucesso.", "id": id})
}

func (h *ProdutoHandler) Listar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}
	produtos, err := h.repo.Listar()
	if err != nil {
		logUnexpected("falha ao listar produtos", err)
		writeError(w, http.StatusInternalServerError, "Não foi possível consultar os produtos")
		return
	}
	writeJSON(w, http.StatusOK, produtos)
}

func (h *ProdutoHandler) BaixarEstoque(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}
	var itens []domain.ItemNota
	if err := decodeJSON(r, &itens); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if len(itens) == 0 {
		writeError(w, http.StatusUnprocessableEntity, "A nota precisa ter pelo menos um item")
		return
	}
	if err := h.repo.BaixarEstoque(itens); err != nil {
		logUnexpected("falha ao baixar estoque", err)
		switch {
		case errors.Is(err, domain.ErrValidacao):
			writeError(w, http.StatusUnprocessableEntity, "Produto e quantidade devem ser válidos")
		case errors.Is(err, domain.ErrProdutoInexistente):
			writeError(w, http.StatusUnprocessableEntity, "Um dos produtos não existe mais")
		case errors.Is(err, domain.ErrSaldoInsuficiente):
			writeError(w, http.StatusUnprocessableEntity, "Saldo insuficiente para concluir a nota")
		default:
			writeError(w, http.StatusInternalServerError, "Não foi possível atualizar o estoque")
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Estoque atualizado"})
}
