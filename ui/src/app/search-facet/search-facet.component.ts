import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { Facet, GameSearch } from '../game.service';

@Component({
  selector: 'app-search-facet',
  templateUrl: './search-facet.component.html',
  styleUrls: ['./search-facet.component.css']
})
export class SearchFacetComponent implements OnInit {
  @Input() facets: Facet[];
  @Input() search: GameSearch;
  @Output() searchEvent = new EventEmitter<GameSearch>();

  constructor() { }

  ngOnInit(): void {
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
