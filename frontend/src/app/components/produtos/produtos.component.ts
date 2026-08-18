import { Component, inject, OnInit } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Produto, ProdutoService } from '../../services/produto.service';

@Component({
  selector: 'app-produtos',
  standalone: true,
  imports: [ReactiveFormsModule],
  templateUrl: './produtos.component.html',
  styleUrl: './produtos.component.css'
})
export class ProdutosComponent implements OnInit {

private produtoService = inject(ProdutoService);

listarProdutos: Produto[] = [];
 
produtoForm = new FormGroup({
  codigo: new FormControl('', [Validators.required]),
  descricao: new FormControl('', [Validators.required]),
  saldo: new FormControl<number | null>(null, [Validators.required, Validators.min(0)])
  });

  ngOnInit() {
    this.carregarProdutos();
  }

  carregarProdutos() {
    this.produtoService.obterProdutos().subscribe({
      next: (produtos) => {
        this.listarProdutos = produtos;
      },
      error: (erro) => {
        console.error('Erro ao obter os produtos:', erro);
      }
    });
  }

onSubmit() {
    if (this.produtoForm.valid) {
      const produtoData = {
        codigo: this.produtoForm.value.codigo!,
        descricao: this.produtoForm.value.descricao!,
        saldo: Number(this.produtoForm.value.saldo)
      };

      this.produtoService.salvarProduto(produtoData).subscribe({
        next: (resposta) => {
          console.log('Produto salvo com sucesso:', resposta);
          alert('Produto salvo com sucesso!');
          this.produtoForm.reset();
          this.carregarProdutos();
        },
        error: (erro) => {
          console.error('Erro ao salvar o produto:', erro);
          alert('Erro ao salvar o produto!');
        }
      });
    } else {
      alert('Por favor, preencha todos os campos corretamente.');
    }
  }
}
