package domain

import "time"

// Produtos que constam na nota
type ItemNota struct {
	ProdutoCodigo string  `json:"produto_codigo"`
	Quantidade    int     `json:"quantidade"`
	PrecoUnitario float64 `json:"preco_unitario"`
	Subtotal      float64 `json:"subtotal"`
}

type NotaFiscal struct {
	ID          int        `json:"id"`
	Numero      int        `json:"numero"`
	DataEmissao time.Time  `json:"data_emissao"`
	ValorTotal  float64    `json:"valor_total"`
	Status      string     `json:"status"`
	Itens       []ItemNota `json:"itens"`
}
