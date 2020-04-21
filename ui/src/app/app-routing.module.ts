import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { GamesComponent } from './games/games.component';
import { FavoritesComponent } from './favorites/favorites.component';

const routes: Routes = [
  { path: 'games', component: GamesComponent },
  { path: 'favorites', component: FavoritesComponent },
  { path: '', redirectTo: '/games', pathMatch: 'full' }
];

@NgModule({
  imports: [RouterModule.forRoot(routes, { useHash: true })],
  exports: [RouterModule]
})
export class AppRoutingModule { }
