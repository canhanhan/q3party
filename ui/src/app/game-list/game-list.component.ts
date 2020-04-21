import { Component, OnInit, Input } from '@angular/core';
import { GameService, Game, Facet, GameSearch, GameResult } from '../game.service';
import { Observable } from 'rxjs';
import { ListGameService } from '../list.service';

interface Service {
  list(GameSearch): Observable<GameResult>;
}

interface GameListModelItem {
  game: Game;
  occupancy: number;
  isFavorite: boolean;
}

interface GameListModel {
  games: GameListModelItem[];
}

@Component({
  selector: 'app-game-list',
  templateUrl: './game-list.component.html',
  styleUrls: ['./game-list.component.css']
})
export class GameListComponent implements OnInit {
  @Input() service: Service;

  facets: Facet[];
  model: GameListModel;
  filters: GameSearch;

  constructor(private favService: ListGameService) {
    this.filters = {};
  }

  ngOnInit(): void {
    this.search(this.filters);
  }

  search(filters: GameSearch): void {
    this.filters = filters;

    this.favService.list(this.filters).subscribe(favs => {
      this.service.list(filters).subscribe((data: GameResult) => {
        this.model = { games: [] };
        data.games.forEach(x => {
          this.model.games.push({
            game: x,
            occupancy: x.clients / x.maxClients * 100,
            isFavorite: favs.games.filter(f => f.server === x.server).length === 1
          });
        });
        this.facets = data.facets;
      });
    });
  }

  addFavorite(game: Game): void {
    const self = this;
    this.favService.add(game).subscribe({
      error(err) { console.error(err); },
      complete() { self.search(self.filters); }
    });
  }

  removeFavorite(game: Game): void {
    const self = this;
    this.favService.remove(game).subscribe({
      error(err) { console.error(err); },
      complete() { self.search(self.filters); }
    });
  }
}
