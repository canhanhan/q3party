import { Component, OnInit } from '@angular/core';
import { GameService, Game, Facet, GameSearch, GameResult } from '../game.service';

@Component({
  selector: 'app-games',
  templateUrl: './games.component.html',
  styleUrls: ['./games.component.css']
})
export class GamesComponent implements OnInit {
  service: GameService;

  constructor(service: GameService) {
    this.service = service;
  }

  ngOnInit(): void {
  }

}
