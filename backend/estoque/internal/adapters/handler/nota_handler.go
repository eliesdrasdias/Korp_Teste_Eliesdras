package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sistema-notas/estoque/internal/core/domain"
	"sistema-notas/estoque/internal/core/ports"
)

type NotaHandler struct {
	repo ports.NotaRepository
}

// NewNotaHandler
func NewNotaHandler(repo ports.NotaRepository) *NotaHandler {
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
		http.Error(w, "Erro interno ao salvar a nota"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mensagem":  "Nota emitida com sucesso",
		"id_gerado": notaID,
	})
}

func (h *NotaHandler) Imprimir(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID int `json:"id"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	nota, err := h.repo.BuscarNotaPorID(req.ID)
	if err != nil {
		http.Error(w, "Nota não encontrada", http.StatusNotFound)
		return
	}

	if nota.Status != "Aberta" {
		http.Error(w, "Apenas notas Abertas podem ser impressas.", http.StatusBadRequest)
		return
	}

	jsonData, _ := json.Marshal(nota.Itens)
	resp, err := http.Post("http://localhost:8080/produtos/baixa", "application/json", bytes.NewBuffer(jsonData))

	// 3. TRATAMENTO DE FALHA OBRIGATÓRIO
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "O Serviço de Estoque está indisponível no momento. A nota não pôde ser impressa.", http.StatusServiceUnavailable)
		return
	}

	h.repo.FecharNota(req.ID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"mensagem": "Nota impressa com sucesso!"})
}
