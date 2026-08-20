import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ItemNota, NotaFiscal, NotaService } from '../services/nota.service';
import { ProdutoService } from '../services/produto.service';
@Component({
  selector: 'app-notas',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './notas.component.html',
  styleUrls: ['./notas.component.css']
})
export class NotasComponent implements OnInit {
  produtosDisponiveis: any[] = [];
  produtoSelecionado: string = '';
  quantidade: number = 1;
  precoUnitario: number = 0;

  itensNota: ItemNota[] = [];
  valorTotalDaNota: number = 0;

  constructor(
    private produtoService: ProdutoService,
    private notaService: NotaService
  ) {}

  ngOnInit() {
    this.carregarProdutos();
  }

  carregarProdutos() {
    this.produtoService.obterProdutos().subscribe({
      next: (dados) => {
        this.produtosDisponiveis = dados;
      },
      error: (err) => console.error('Erro ao carregar produtos', err)
    });
  }

  // Função que adiciona um item na nota
  adicionarItem() {
    if (!this.produtoSelecionado || this.quantidade <= 0 || this.precoUnitario <= 0) {
      alert("Preencha o produto, quantidade e preço corretamente.");
      return;
    }

    const subtotal = this.quantidade * this.precoUnitario;
    
    this.itensNota.push({
      produto_codigo: this.produtoSelecionado,
      quantidade: this.quantidade,
      preco_unitario: this.precoUnitario,
      subtotal: subtotal
    });

    this.valorTotalDaNota += subtotal;
    this.produtoSelecionado = '';
    this.quantidade = 1;
    this.precoUnitario = 0;
  }

  // Função para emitir a nota
  emitirNota() {
    if (this.itensNota.length === 0) {
      alert("Sua nota precisa ter pelo menos um item.");
      return;
    }

    const nota: NotaFiscal = {
      valor_total: this.valorTotalDaNota,
      itens: this.itensNota
    };

    this.notaService.emitirNota(nota).subscribe({
      next: (resposta) => {
        alert(`Sucesso! ${resposta.mensagem} (ID: ${resposta.id_gerado})`);
        this.itensNota = [];
        this.valorTotalDaNota = 0;
      },
      error: (err) => {
        alert("Erro ao emitir a nota fiscal. Verifique o console.");
        console.error(err);
      }
    });
  }
}