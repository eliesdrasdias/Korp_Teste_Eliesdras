import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export interface ItemNota {
  produto_codigo: string;
  quantidade: number;
  preco_unitario: number;
  subtotal: number;
}

export interface NotaFiscal {
  valor_total: number;
  itens: ItemNota[];
}

@Injectable({
  providedIn: 'root'
})
export class NotaService {
  private apiUrl = 'http://localhost:8080/notas';

  constructor(private http: HttpClient) { }

// Método para emitir uma nota
  emitirNota(nota: NotaFiscal): Observable<any> {
    return this.http.post(this.apiUrl, nota);
  }
}