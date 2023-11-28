import {Injectable, Injector} from '@angular/core';
import {BaseService} from "../../core/services/base.service";
import {KanbanColumn} from "../models/kanban-column";

@Injectable({
  providedIn: 'root'
})
export class KanbanColumnService extends BaseService<KanbanColumn>{
  public override usePagination = false

  constructor(protected override injector: Injector) {
    super(injector)
    this.apiUrl = this.config.apiUrl + '/columns'
  }
}
