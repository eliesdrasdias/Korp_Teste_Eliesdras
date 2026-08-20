package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sistema-notas/estoque/internal/core/domain"
	"sistema-notas/estoque/internal/core/ports"
)

type ProdutoHandler struct {
	repo ports.ProdutoRepository
}

// NewProdutoHandler
func NewProdutoHandler(repo ports.ProdutoRepository) *ProdutoHandler {
	return &ProdutoHandler{repo: repo}
}

// Criar
func (h *ProdutoHandler) Criar(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var p domain.Produto
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Erro ao decodificar o JSON", http.StatusBadRequest)
		return
	}

	id, err := h.repo.Salvar(p)
	if err != nil {
		fmt.Println("Falha no banco de dados", err)
		http.Error(w, "Erro interno ao salvar o produto", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"mensagem": "Produto salvo com sucesso", "id": %d}`, id)
}

func (h *ProdutoHandler) Listar(w http.ResponseWriter, r *http.Request) {
	produtos, err := h.repo.Listar()
	if err != nil {
		fmt.Println("Falha no banco de dados", err)
		http.Error(w, "Erro interno ao listar os produtos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(produtos)
}

func (h *ProdutoHandler) BaixarEstoque(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var itens []domain.ItemNota
	if err := json.NewDecoder(r.Body).Decode(&itens); err != nil {
		http.Error(w, "Erro ao ler os dados enviados", http.StatusBadRequest)
		return
	}

	if err := h.repo.BaixarEstoque(itens); err != nil {
		http.Error(w, "Erro interno ao atualizar estoque", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
