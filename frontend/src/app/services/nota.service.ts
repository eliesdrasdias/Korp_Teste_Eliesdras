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
	  id?: number;
	  numero?: number;
	  status?: 'ABERTA' | 'FECHADA';
  valor_total: number;
  itens: ItemNota[];
}

@Injectable({
  providedIn: 'root'
})
export class NotaService {
  private apiUrl = 'http://localhost:8081/notas';

  constructor(private http: HttpClient) { }

// Método para emitir uma nota
  emitirNota(nota: NotaFiscal): Observable<NotaFiscal> {
    return this.http.post<NotaFiscal>(this.apiUrl, nota);
  }

  listarNotas(): Observable<NotaFiscal[]> { return this.http.get<NotaFiscal[]>(`${this.apiUrl}/listar`); }

  imprimirNota(id: number): Observable<{ message: string; status: 'FECHADA' }> {
    return this.http.post<{ message: string; status: 'FECHADA' }>('http://localhost:8081/notas/imprimir', { id });
  }
}
