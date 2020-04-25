import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from './../environments/environment';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { Game, GameSearch, GameResult, GameService, filter } from './game.service';

interface ListResponse {
  id: string;
  servers: string[];
}

@Injectable({
  providedIn: 'root'
})

export class ListGameService {
  private id: string;

  constructor(private listService: ListService, private gameService: GameService) {
    this.id = '5E965D56';
  }

  list(s: GameSearch): Observable<GameResult> {
    return new Observable<GameResult>(observer => {
      this.gameService.list(s).subscribe(res => {
        this.listService.get(this.id).subscribe(list => {
          const listGames = list.servers?.map(v => { 
            const r = res.games.filter(g => g.server === v);
            return r.length < 1 ? null : r[0];
          }).filter(g => g != null);
          observer.next(filter(listGames ?? [], s));
          observer.complete();
        });
      });
    });
  }

  add(game: Game) {
    return this.listService.addToList(this.id, game.server);
  }

  remove(game: Game) {
    return this.listService.removeFromList(this.id, game.server);
  }
}

@Injectable({
  providedIn: 'root'
})
export class ListService {
  id: string;

  constructor(private http: HttpClient) { }

  create(): Observable<ListResponse> {
    return this.http.post<ListResponse>(environment.apiUrl + '/lists', null);
  }

  get(id: string): Observable<ListResponse> {
    return this.http.get<ListResponse>(environment.apiUrl + '/lists/' + id);
  }

  addToList(id: string, server: string): Observable<ListResponse> {
    return this.http.post<ListResponse>(environment.apiUrl + '/lists/' + id + '/add', {server});
  }

  removeFromList(id: string, server: string): Observable<ListResponse> {
    return this.http.post<ListResponse>(environment.apiUrl + '/lists/' + id + '/remove', {server});
  }
}
