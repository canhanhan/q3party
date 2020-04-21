import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { finalize, tap, map } from 'rxjs/operators';
import { environment } from '../environments/environment';

export interface FacetItem {
  name: string;
  count: number;
  selected: boolean;
}

export interface Facet {
  name: string;
  title: string;
  values: FacetItem[];
}

export interface Game {
  server: string;

  clients: number;
  maxClients: number;
  humanPlayers: number;

  needPassword: boolean;
  isPure: boolean;

  game: string;
  map: string;
  name: string;
}

export interface GameSearch {
  map?: string;
  game?: string;
  needPassword?: string;
  isPure?: string;
}

export interface GameResult {
  facets: Facet[];
  games: Game[];
}

export interface GameListItem {
  server: string;
  info: {[key: string]: any};
  status: {[key: string]: any};
}

export type facetValue = {[key: string]: (x: any) => boolean};

@Injectable({
  providedIn: 'root'
})
export class GameService {
  constructor(private http: HttpClient) {}

  list(s: GameSearch): Observable<GameResult> {
    return this.http.get<GameListItem[]>(environment.apiUrl + '/games').pipe(map((data, i) => {
      return filter(data.map(v => {
        return {
          server: v.server,
          clients: parseInt(v.info.clients ?? '0', 10),
          maxClients: parseInt(v.info.sv_maxclients ?? '0', 10),
          humanPlayers: parseInt(v.info.g_humanplayers ?? '0', 10),
          needPassword: v.info.g_needpass === '1',
          isPure: v.info.pure === '1',
          game: v.info.game,
          map: v.info.mapname,
          name: v.info.hostname
        };
      }), s);
    }));
  }
}

export function filter(data: Game[], s: GameSearch): GameResult {
  const games = data.filter(x => {
    if (s && s.map != null && x.map !== s.map) {
      return false;
    }

    if (s && s.game != null && ((s.game === 'No Mods' && x.game != null) || (s.game !== 'No Mods' && x.game !== s.game))) {
      return false;
    }

    if (s && s.needPassword != null && x.needPassword !== (s.needPassword === 'Yes')) {
      return false;
    }

    if (s && s.isPure != null && x.isPure !== (s.isPure === 'Yes')) {
      return false;
    }

    return true;
  });

  const facets = makeFacets(games);

  return {
    games,
    facets
  };
}

function makeFacets(games: Game[]): Facet[] {
  const facets = [];
  facets.push(calculateFacet(games, 'Maps', 'map'));
  facets.push(calculateFacet(games, 'Game', 'game', null, {'No Mods': x => !x}));
  facets.push(calculateFacet(games, 'Pure', 'isPure', {Yes: x => x, No: x => !x}));
  facets.push(calculateFacet(games, 'Password', 'needPassword', {Yes: x => x, No: x => !x}));

  return facets;
}

function calculateFacet(games: Game[], name: string, prop: string, values?: facetValue, customValues?: facetValue): Facet {
  if (!values) {
    values = {};
    const data = [...new Set(games.map((x: Game) => x[prop]))];
    for (const value of data.filter(x => x != null)) {
      values[value] = x => x === value;
    }
  }

  if (customValues) {
    Object.assign(values, customValues);
  }

  const facet: Facet = { name: prop, title: name, values: [] };

  Object.getOwnPropertyNames(values).forEach(title => {
    const matcher = values[title];
    const cnt = games.filter(g => matcher(g[prop])).length;
    if (cnt > 0) {
      facet.values.push({
        name: title,
        count: cnt,
        selected: false,
      });
    }
  });

  facet.values.sort((a, b) => {
    if (a.count > b.count) {
      return -1;
    } else if (a.count < b.count) {
      return 1;
    } else {
      return 0;
    }
  });

  return facet;
}
