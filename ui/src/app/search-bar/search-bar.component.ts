import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { Facet, GameSearch } from '../game.service';

interface playerFilter {
  min: number;
  max: number;
}

@Component({
  selector: 'app-search-bar',
  templateUrl: './search-bar.component.html',
  styleUrls: ['./search-bar.component.css']
})
export class SearchBarComponent implements OnInit {
  @Input() facets: Facet[];
  @Input() search: GameSearch;
  @Output() searchEvent = new EventEmitter<GameSearch>();

  constructor() { }

  ngOnInit(): void {
  }

  update(name: string, value: string) {
    const val = isNaN(parseInt(value, 10)) ? null : parseInt(value, 10);
    switch (name) {
      case 'minPlayers':
        this.search.minPlayers = val;
        break;
      case 'maxPlayers':
        this.search.maxPlayers = val;
        break;
      case 'minBots':
        this.search.minBots = val;
        break;
      case 'maxBots':
        this.search.maxBots = val;
        break;
      case 'minPing':
        this.search.minPing = val;
        break;
      case 'maxPing':
        this.search.maxPing = val;
        break;
      default:
        throw new Error('Unknown field: ' + name);
    }

    this.searchEvent.emit(this.search);
  }

  click(facet: string, value: any) {
    switch (facet) {
      case 'map':
        this.search.map = this.search.map === value ? null : value;
        break;
      case 'game':
        this.search.game = this.search.game === value ? null : value;
        break;
      case 'needPassword':
        this.search.needPassword = this.search.needPassword === value ? null : value;
        break;
      case 'isPure':
        this.search.isPure = this.search.isPure === value ? null : value;
        break;
      default:
        return;
    }

    this.searchEvent.emit(this.search);
  }
}
