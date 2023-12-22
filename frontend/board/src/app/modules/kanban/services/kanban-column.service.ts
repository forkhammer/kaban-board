import {Injectable, Injector} from '@angular/core';
import {BaseService} from "../../core/services/base.service";
import {KanbanColumn} from "../models/kanban-column";
import { map } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class KanbanColumnService extends BaseService<KanbanColumn>{
  public override usePagination = false

  constructor(protected override injector: Injector) {
    super(injector)
    this.apiUrl = this.config.apiUrl + '/columns'
  }

  saveOrdering(columns: KanbanColumn[]) {
    const data = columns.map((column, index) => {
      return {
        id: column.id,
        order: column.order
      }
    })
    return this.http.post(`${this.apiUrl}/save_ordering`, data).pipe(
      map(items => items as KanbanColumn[])
    )
  }
}
