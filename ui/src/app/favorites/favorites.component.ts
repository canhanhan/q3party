import { Component, OnInit } from '@angular/core';
import { ListGameService } from '../list.service';

@Component({
  selector: 'app-favorites',
  templateUrl: './favorites.component.html',
  styleUrls: ['./favorites.component.css']
})
export class FavoritesComponent implements OnInit {
  service: ListGameService;

  constructor(service: ListGameService) {
    this.service = service;
  }

  ngOnInit(): void {
  }

}
