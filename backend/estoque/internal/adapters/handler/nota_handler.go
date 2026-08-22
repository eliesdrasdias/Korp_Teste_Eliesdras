package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"sistema-notas/estoque/internal/core/domain"
	"sistema-notas/estoque/internal/core/ports"
)

type NotaHandler struct {
	repo     ports.NotaRepository
	stockURL string
	client   *http.Client
}

func NewNotaHandler(repo ports.NotaRepository, stockURL string) *NotaHandler {
	return &NotaHandler{repo: repo, stockURL: strings.TrimRight(stockURL, "/"), client: &http.Client{Timeout: 4 * time.Second}}
}

func (h *NotaHandler) Emitir(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}
	var nota domain.NotaFiscal
	if err := decodeJSON(r, &nota); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if err := validarNota(&nota); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	created, err := h.repo.Emitir(nota)
	if err != nil {
		logUnexpected("falha ao emitir nota", err)
		writeError(w, http.StatusInternalServerError, "Não foi possível criar a nota")
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *NotaHandler) Listar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}
	notas, err := h.repo.Listar()
	if err != nil {
		logUnexpected("falha ao listar notas", err)
		writeError(w, http.StatusInternalServerError, "Não foi possível consultar as notas")
		return
	}
	writeJSON(w, http.StatusOK, notas)
}

func (h *NotaHandler) Imprimir(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}
	var req struct {
		ID int `json:"id"`
	}
	if err := decodeJSON(r, &req); err != nil || req.ID <= 0 {
		writeError(w, http.StatusBadRequest, "Identificador da nota inválido")
		return
	}
	nota, err := h.repo.BuscarNotaPorID(req.ID)
	if err != nil {
		writeError(w, http.StatusNotFound, "Nota não encontrada")
		return
	}
	if nota.Status == "FECHADA" {
		writeError(w, http.StatusConflict, "Esta nota já está fechada e não pode ser impressa novamente")
		return
	}
	if len(nota.Itens) == 0 {
		writeError(w, http.StatusUnprocessableEntity, "Não é possível fechar uma nota sem itens")
		return
	}
	log.Printf("início do fechamento da nota %d", nota.ID)
	payload, _ := json.Marshal(nota.Itens)
	resp, err := h.client.Post(h.stockURL+"/produtos/baixa", "application/json", bytes.NewReader(payload))
	if err != nil {
		logUnexpected("serviço de estoque indisponível", err)
		writeError(w, http.StatusServiceUnavailable, "Não foi possível concluir a nota porque o serviço de estoque está temporariamente indisponível. Tente novamente em alguns instantes.")
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnprocessableEntity {
		writeError(w, http.StatusUnprocessableEntity, "Não foi possível concluir a nota: estoque insuficiente ou produto indisponível")
		return
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		logUnexpected("estoque retornou resposta inesperada: "+string(body), errors.New(resp.Status))
		writeError(w, http.StatusServiceUnavailable, "Não foi possível concluir a nota porque o serviço de estoque está temporariamente indisponível. Tente novamente em alguns instantes.")
		return
	}
	if err := h.repo.FecharNota(nota.ID); err != nil {
		logUnexpected("falha ao fechar nota após baixa de estoque", err)
		if errors.Is(err, domain.ErrNotaFechada) {
			writeError(w, http.StatusConflict, "Esta nota já está fechada")
			return
		}
		writeError(w, http.StatusInternalServerError, "A baixa foi confirmada, mas não foi possível finalizar a nota. Consulte o suporte.")
		return
	}
	log.Printf("nota %d fechada com sucesso", nota.ID)
	writeJSON(w, http.StatusOK, map[string]any{"message": "Nota fechada com sucesso.", "nota": nota.ID, "numero": nota.Numero, "status": "FECHADA"})
}

func validarNota(nota *domain.NotaFiscal) error {
	if len(nota.Itens) == 0 {
		return domain.ErrNotaSemItens
	}
	var total float64
	for i := range nota.Itens {
		item := &nota.Itens[i]
		item.ProdutoCodigo = strings.TrimSpace(item.ProdutoCodigo)
		if item.ProdutoCodigo == "" || item.Quantidade <= 0 || item.PrecoUnitario < 0 {
			return domain.ErrValidacao
		}
		item.Subtotal = float64(item.Quantidade) * item.PrecoUnitario
		total += item.Subtotal
	}
	nota.ValorTotal = total
	return nil
}
