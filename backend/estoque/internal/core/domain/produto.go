package domain

type Produto struct {
	ID        int     `json:"id"`
	Codigo    string  `json:"codigo"`
	Descricao string  `json:"descricao"`
	Saldo     float64 `json:"saldo"`
}
