import { Routes } from '@angular/router';
import { ProdutosComponent } from './components/produtos/produtos.component';

export const routes: Routes = [
    { path: 'produtos', component: ProdutosComponent },
    { path: '', redirectTo: '/produtos', pathMatch: 'full' }
];
