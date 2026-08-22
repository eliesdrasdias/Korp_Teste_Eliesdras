import { Component, DestroyRef, inject, OnInit } from '@angular/core';
import { ReactiveFormsModule, FormControl, FormGroup, Validators } from '@angular/forms';
import { finalize } from 'rxjs';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { Produto, ProdutoService } from '../../services/produto.service';
@Component({selector:'app-produtos',standalone:true,imports:[ReactiveFormsModule],templateUrl:'./produtos.component.html',styleUrl:'./produtos.component.css'})
export class ProdutosComponent implements OnInit {
  private service=inject(ProdutoService); private destroyRef=inject(DestroyRef); listarProdutos:Produto[]=[]; carregando=false; salvando=false; erro=''; mensagem='';
  produtoForm=new FormGroup({codigo:new FormControl('',{nonNullable:true,validators:[Validators.required]}),descricao:new FormControl('',{nonNullable:true,validators:[Validators.required]}),saldo:new FormControl<number|null>(null,[Validators.required,Validators.min(0)])});
  ngOnInit():void { this.carregarProdutos(); }
  carregarProdutos():void { this.carregando=true; this.erro=''; this.service.obterProdutos().pipe(takeUntilDestroyed(this.destroyRef),finalize(()=>this.carregando=false)).subscribe({next:p=>this.listarProdutos=p,error:()=>this.erro='Não foi possível carregar os produtos. Verifique o serviço de estoque.'}); }
  onSubmit():void { this.mensagem=''; this.erro=''; if(this.produtoForm.invalid){this.produtoForm.markAllAsTouched();return;} this.salvando=true; this.service.salvarProduto(this.produtoForm.getRawValue() as Produto).pipe(takeUntilDestroyed(this.destroyRef),finalize(()=>this.salvando=false)).subscribe({next:()=>{this.mensagem='Produto cadastrado com sucesso.';this.produtoForm.reset();this.carregarProdutos();},error:e=>this.erro=e.error?.message??'Não foi possível cadastrar o produto.'}); }
}
