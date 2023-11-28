import { Injectable } from '@angular/core';
import {BehaviorSubject} from "rxjs";

@Injectable({
  providedIn: 'root'
})
export class ThemeServiceService {
  theme$ = new BehaviorSubject<string>('light')

  constructor() {
    this.loadSettings()

    this.theme$.subscribe(value => {
      this.saveSettings()
    })
  }

  loadSettings() {
    const theme = localStorage.getItem('dark-mode') === '1' ? 'dark': 'light'
    this.theme$.next(theme);
  }

  saveSettings() {
    localStorage.setItem('dark-mode', this.theme$.value === 'dark' ? '1' : '0')
  }
}
