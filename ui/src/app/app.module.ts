import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { NavItemComponent } from './nav-item/nav-item.component';
import { GameListComponent } from './game-list/game-list.component';
import { Q3ColorPipe } from './q3color.pipe';
import { SearchFacetComponent } from './search-facet/search-facet.component';
import { MapNamePipe } from './map-name.pipe';
import { GamesComponent } from './games/games.component';
import { FavoritesComponent } from './favorites/favorites.component';

@NgModule({
  declarations: [
    AppComponent,
    NavItemComponent,
    GameListComponent,
    Q3ColorPipe,
    SearchFacetComponent,
    MapNamePipe,
    GamesComponent,
    FavoritesComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule
  ],
  providers: [
    { provide: Window, useValue: window }
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
