import { CommonModule } from '@angular/common';
import { Component, DestroyRef, inject, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { finalize } from 'rxjs';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { ItemNota, NotaFiscal, NotaService } from '../services/nota.service';
import { Produto, ProdutoService } from '../services/produto.service';

interface ItemNotaTela extends ItemNota {
  descricao: string;
}

@Component({
  selector: 'app-notas',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './notas.component.html',
  styleUrls: ['./notas.component.css']
})
export class NotasComponent implements OnInit {
  private readonly destroyRef = inject(DestroyRef);
  private readonly produtoService = inject(ProdutoService);
  private readonly notaService = inject(NotaService);

  produtosDisponiveis: Produto[] = [];
  notas: NotaFiscal[] = [];
  itensNota: ItemNotaTela[] = [];

  produtoSelecionado = '';
  quantidade = 1;
  precoUnitario = 0;
  notaGerada: NotaFiscal | null = null;

  carregandoProdutos = false;
  carregandoNotas = false;
  emitindo = false;
  imprimindo = false;
  erro = '';
  mensagem = '';

  get valorTotalDaNota(): number {
    return (this.itensNota ?? []).reduce((total, item) => total + item.subtotal, 0);
  }

  get produtoAtual(): Produto | undefined {
    return (this.produtosDisponiveis ?? []).find(produto => produto.codigo === this.produtoSelecionado);
  }

  ngOnInit(): void {
    this.carregarProdutos();
    this.carregarNotas();
  }

  carregarProdutos(): void {
    this.carregandoProdutos = true;

    this.produtoService.obterProdutos()
      .pipe(
        takeUntilDestroyed(this.destroyRef),
        finalize(() => this.carregandoProdutos = false)
      )
      .subscribe({
        next: produtos => this.produtosDisponiveis = Array.isArray(produtos) ? produtos : [],
        error: () => this.erro = 'Não foi possível carregar os produtos do serviço de estoque.'
      });
  }

  carregarNotas(): void {
    this.carregandoNotas = true;

    this.notaService.listarNotas()
      .pipe(
        takeUntilDestroyed(this.destroyRef),
        finalize(() => this.carregandoNotas = false)
      )
      .subscribe({
        next: notas => this.notas = Array.isArray(notas) ? notas : [],
        // A indisponibilidade da listagem de notas não bloqueia a seleção de produtos.
        error: () => this.erro = 'Não foi possível carregar as notas já emitidas.'
      });
  }

  adicionarItem(): void {
    this.erro = '';
    const produto = this.produtoAtual;

    if (!produto) {
      this.erro = 'Selecione um produto.';
      return;
    }

    if (!Number.isInteger(this.quantidade) || this.quantidade <= 0 || this.precoUnitario <= 0) {
      this.erro = 'Informe uma quantidade inteira maior que zero e um preço válido.';
      return;
    }

    const quantidadeJaAdicionada = (this.itensNota ?? [])
      .filter(item => item.produto_codigo === produto.codigo)
      .reduce((total, item) => total + item.quantidade, 0);

    if (quantidadeJaAdicionada + this.quantidade > produto.saldo) {
      this.erro = `Quantidade maior que o saldo disponível de ${produto.descricao} (${produto.saldo}).`;
      return;
    }

    // Mantém a lista sempre utilizável mesmo diante de uma alteração indevida de estado.
    this.itensNota = this.itensNota ?? [];
    this.itensNota.push({
      produto_codigo: produto.codigo,
      descricao: produto.descricao,
      quantidade: this.quantidade,
      preco_unitario: this.precoUnitario,
      subtotal: this.quantidade * this.precoUnitario
    });

    this.produtoSelecionado = '';
    this.quantidade = 1;
    this.precoUnitario = 0;
  }

  removerItem(index: number): void {
    this.itensNota.splice(index, 1);
  }

  emitirNota(): void {
    if (!this.itensNota?.length) {
      this.erro = 'Adicione pelo menos um produto à nota.';
      return;
    }

    this.emitindo = true;
    this.erro = '';
    this.mensagem = '';

    // A descrição é usada apenas na interface; o contrato do backend recebe código e valores.
    const itens: ItemNota[] = (this.itensNota ?? []).map(({ descricao, ...item }) => item);
    const nota: NotaFiscal = { valor_total: this.valorTotalDaNota, itens };

    this.notaService.emitirNota(nota)
      .pipe(takeUntilDestroyed(this.destroyRef), finalize(() => this.emitindo = false))
      .subscribe({
        next: notaCriada => {
          this.notaGerada = notaCriada;
          this.itensNota = [];
          this.mensagem = `Nota ${notaCriada.numero} criada como ABERTA.`;
          this.carregarNotas();
        },
        error: error => this.erro = error.error?.message ?? 'Não foi possível criar a nota.'
      });
  }

  imprimirNota(): void {
    if (!this.notaGerada || this.notaGerada.status === 'FECHADA' || this.imprimindo) return;

    this.imprimindo = true;
    this.erro = '';
    this.mensagem = '';

    this.notaService.imprimirNota(this.notaGerada.id!)
      .pipe(takeUntilDestroyed(this.destroyRef), finalize(() => this.imprimindo = false))
      .subscribe({
        next: resposta => {
          this.notaGerada = { ...this.notaGerada!, status: resposta.status };
          this.mensagem = resposta.message;
          this.carregarProdutos();
          this.carregarNotas();
        },
        error: error => this.erro = error.error?.message ?? 'Não foi possível fechar a nota.'
      });
  }
}
