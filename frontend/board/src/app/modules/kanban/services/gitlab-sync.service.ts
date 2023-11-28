import { Injectable } from '@angular/core';
import {BehaviorSubject} from "rxjs";

@Injectable({
  providedIn: 'root'
})
export class GitlabSyncService {
  public updateTime$ = new BehaviorSubject<Date | null>(null)

  constructor() { }
}
