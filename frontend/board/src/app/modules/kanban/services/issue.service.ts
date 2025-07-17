import {Injectable, Injector} from '@angular/core';
import {BaseService} from "../../core/services/base.service";
import { KanbanIssue } from '../models/kanban-issue';

@Injectable({
  providedIn: 'root'
})
export class IssueService extends BaseService<KanbanIssue>{
  public override usePagination = false
  override useCache = false

  constructor(protected override injector: Injector) {
    super(injector)
    this.apiUrl = this.config.apiUrl + '/issue'
  }

  formatter(item: KanbanIssue) {
    return `#${item.projectName}${item.iid} - ${item.title}`
  }

  bindToSprint(issueId: string, sprintId: number) {
    return this.http.post<KanbanIssue>(`${this.apiUrl}/${issueId}/bind`, {sprint_id: sprintId})
  }

  unbindFromSprint(issueId: string, bindingId: number) {
    return this.http.post<KanbanIssue>(`${this.apiUrl}/${issueId}/unbind`, {binding_id: bindingId})
  }
}
