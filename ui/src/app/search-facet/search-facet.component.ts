import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'app-search-facet',
  templateUrl: './search-facet.component.html',
  styleUrls: ['./search-facet.component.css']
})
export class SearchFacetComponent implements OnInit {
  @Input() name: string;
  @Input() title: string;

  constructor() { }

  ngOnInit(): void {
  }

}
