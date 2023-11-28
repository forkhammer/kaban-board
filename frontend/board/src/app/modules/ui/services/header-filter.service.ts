import { Injectable } from '@angular/core';
import {BehaviorSubject} from "rxjs";
import {Team} from "../../kanban/models/team";

@Injectable({
  providedIn: 'root'
})
export class HeaderFilterService {
  public team$ = new BehaviorSubject<Team | null>(null);

  constructor() { }
}
