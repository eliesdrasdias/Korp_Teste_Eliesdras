import { Component, inject } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { ProdutoService } from '../../services/produto.service';

@Component({
  selector: 'app-produtos',
  standalone: true,
  imports: [ReactiveFormsModule],
  templateUrl: './produtos.component.html',
  styleUrl: './produtos.component.css'
})
export class ProdutosComponent {

private produtoService = inject(ProdutoService);
 
produtoForm = new FormGroup({
  codigo: new FormControl('', [Validators.required]),
  descricao: new FormControl('', [Validators.required]),
  saldo: new FormControl<number | null>(null, [Validators.required, Validators.min(0)])
  });

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
