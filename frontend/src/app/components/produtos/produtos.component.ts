import { Component } from '@angular/core';
import { ReactiveFormsModule, FormGroup, FormControl, Validators } from '@angular/forms';

@Component({
  selector: 'app-produtos',
  standalone: true,
  imports: [ReactiveFormsModule],
  templateUrl: './produtos.component.html',
  styleUrl: './produtos.component.css'
})
export class ProdutosComponent {
  produtoForm = new FormGroup({
    codigo: new FormControl('', [Validators.required]),
    descricao: new FormControl('', [Validators.required]),
    saldo: new FormControl<number | null>(null, [Validators.required, Validators.min(0)])
  });

  onSubmit() {
    if (this.produtoForm.valid) {
      console.log('Produto pronto para ser salvo', this.produtoForm.value);
      alert('Produto capturado com sucesso! Olhe o console (F12)');
      this.produtoForm.reset();
    } else {
      alert('Por favor, preencha todos os campos corretamente.');
    }
  }
}
