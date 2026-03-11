import {Injectable, Injector} from '@angular/core';
import {BaseService} from "../../core/services/base.service";
import { KanbanIssue } from '../models/kanban-issue';
import { map, Observable, tap } from 'rxjs';

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

  saveOrdering(issues: KanbanIssue[]): Observable<void> {
    const data = issues
      .filter(issue => issue.bindingId != null)
      .map(issue => ({
        id: issue.bindingId,
        order: issue.order ?? '',
      }));
    return this.http.post<void>(`${this.apiUrl}/save_ordering`, data);
  }
}
