import { inject, Injectable } from '@angular/core';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root',
})
export class NotFoundService {
  private router = inject(Router)

  setNotFound() {
    this.router.navigate(['not-found']);
  }
}
