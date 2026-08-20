import { Routes } from '@angular/router';
import { ProdutosComponent } from './components/produtos/produtos.component';
import { NotasComponent } from './notas/notas.component';

export const routes: Routes = [
    { path: 'produtos', component: ProdutosComponent },
    { path: 'notas', component: NotasComponent },
    { path: '', redirectTo: '/produtos', pathMatch: 'full' }
];
