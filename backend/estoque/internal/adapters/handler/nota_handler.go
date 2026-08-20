package handler

import (
	"encoding/json"
	"net/http"
	"sistema-notas/estoque/internal/core/domain"
	"sistema-notas/estoque/internal/core/ports"
)

type NotaHandler struct {
	repo ports.NotaRepository
}

// NewNotaHandler
func NovoNotaHandler(repo ports.NotaRepository) *NotaHandler {
	return &NotaHandler{repo: repo}
}

// Emitir
func (h *NotaHandler) Emitir(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var nota domain.NotaFiscal
	err := json.NewDecoder(r.Body).Decode(&nota)
	if err != nil {
		http.Error(w, "Erro ao ler os dados enviados", http.StatusBadRequest)
		return
	}

	if len(nota.Itens) == 0 {
		http.Error(w, "A nota fiscal precisa ter pelo menos um item", http.StatusBadRequest)
		return
	}

	notaID, err := h.repo.Emitir(nota)
	if err != nil {
		http.Error(w, "Erro interno ao salvar a nota", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mensagem":  "Nota emitida com sucesso",
		"id_gerado": notaID,
	})
}
