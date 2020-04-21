import { Component, OnInit, Input } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-nav-item',
  templateUrl: './nav-item.component.html',
  styleUrls: ['./nav-item.component.css']
})

export class NavItemComponent {
  @Input() route: string;
  @Input() title: string;

  constructor(private router: Router) { }

  currentRoute(): string {
    return this.router.url;
  }
}
