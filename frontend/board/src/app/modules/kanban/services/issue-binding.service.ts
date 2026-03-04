import {Injectable, Injector} from '@angular/core';
import {BaseService} from "../../core/services/base.service";
import { IssueBinding, KanbanIssue } from '../models/kanban-issue';
import { map, tap } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class IssueBindingService extends BaseService<KanbanIssue>{
  public override usePagination = true
  override useCache = false

  constructor(protected override injector: Injector) {
    super(injector)
    this.apiUrl = this.config.apiUrl + '/binding'
  }

  formatter(item: IssueBinding) {
    return `#${item.projectName}${item.iid} - ${item.title}`
  }

  override delete(data: KanbanIssue) {
    return this.http.delete(`${this.apiUrl}/${data.bindingId}`).pipe(
      map(res => res as KanbanIssue),
      tap(this.updateItemCache.bind(this)),
    );
  }

  override save(data: any) {
    return this.http.put(`${this.apiUrl}/${data.bindingId}`, data).pipe(
      map(res => res as KanbanIssue),
      tap(this.updateItemCache.bind(this)),
    );
  }
}
