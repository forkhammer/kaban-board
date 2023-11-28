import {Injectable, Injector} from '@angular/core';
import {BaseService} from "../../core/services/base.service";
import {KanbanColumn} from "../models/kanban-column";
import {Team} from "../models/team";

@Injectable({
  providedIn: 'root'
})
export class TeamService extends BaseService<Team>{
  public override usePagination = false

  constructor(protected override injector: Injector) {
    super(injector)
    this.apiUrl = this.config.apiUrl + '/teams'
  }
}
