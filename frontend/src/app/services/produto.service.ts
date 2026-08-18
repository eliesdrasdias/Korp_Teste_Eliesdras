import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export interface Produto {
  codigo: string;
  descricao: string;
  saldo: number;
}

@Injectable({
  providedIn: 'root'
})
export class ProdutoService {
  private http = inject(HttpClient);
  private apiUrl = 'http://localhost:8080/produtos';

  // Método para salvar um produto
  salvarProduto(produto: Produto): Observable<any> {
    return this.http.post(this.apiUrl, produto);
  }

  // Método para obter todos os produtos
  obterProdutos(): Observable<Produto[]> {
    return this.http.get<Produto[]>(`${this.apiUrl}/listar`);
  }
}
