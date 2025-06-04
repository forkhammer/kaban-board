import {Injectable, Injector} from '@angular/core';
import {BaseService} from "../../core/services/base.service";
import { Sprint } from '../models/sprint';

@Injectable({
  providedIn: 'root'
})
export class SprintService extends BaseService<Sprint>{
  public override usePagination = false

  constructor(protected override injector: Injector) {
    super(injector)
    this.apiUrl = this.config.apiUrl + '/sprint'
  }

  complete(sprintId: number) {
    return this.http.post<Sprint>(`${this.apiUrl}/${sprintId}/complete`, {})
  }

  run(sprintId: number) {
    return this.http.post<Sprint>(`${this.apiUrl}/${sprintId}/run`, {})
  }
}
