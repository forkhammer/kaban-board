import {Injectable, Injector} from '@angular/core';
import {BaseService} from "../../core/services/base.service";
import {KanbanUser, KanbanUserResponse} from "../models/kanban-user";
import {RestQuery} from "../../core/models/rest";
import {map} from "rxjs/operators";
import {Pagination} from "../../core/models/base";

@Injectable({
  providedIn: 'root'
})
export class KanbanUserService extends BaseService<KanbanUser>{
  public override usePagination = false

  constructor(protected override injector: Injector) {
    super(injector)
    this.apiUrl = this.config.apiUrl + '/kanban-users'
  }

  listUsers(query?: RestQuery) {
    return this.http.get(this.apiUrl, { params: this.filterQuery(query) }).pipe(
      map(res => res as KanbanUserResponse)
    );
  }
}
