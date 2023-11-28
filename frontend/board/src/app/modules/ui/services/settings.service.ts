import { Injectable } from '@angular/core';
import {BehaviorSubject} from "rxjs";

@Injectable({
  providedIn: 'root'
})
export class SettingsService {
  public isOpen$ = new BehaviorSubject<boolean>(false)

  constructor() { }

  open() {
    this.isOpen$.next(true)
  }

  close() {
    this.isOpen$.next(false)
  }

  toggle() {
    this.isOpen$.next(!this.isOpen$.value)
  }
}
