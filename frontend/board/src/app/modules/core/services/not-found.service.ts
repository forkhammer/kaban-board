import { Injectable } from '@angular/core';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root',
})
export class NotFoundService {
  constructor(private router: Router) {}

  setNotFound() {
    this.router.navigate(['not-found']);
  }
}
