import {Injectable, Injector} from '@angular/core';
import {BaseService} from "../../core/services/base.service";
import { KanbanIssue } from '../models/kanban-issue';

@Injectable({
  providedIn: 'root'
})
export class IssueService extends BaseService<KanbanIssue>{
  public override usePagination = true
  override useCache = false

  constructor(protected override injector: Injector) {
    super(injector)
    this.apiUrl = this.config.apiUrl + '/issue'
  }

  formatter(item: KanbanIssue) {
    return `#${item.projectName}${item.iid} - ${item.title}`
  }

  bindToSprint(issueIds: number[], sprintId: number, assigneeId: number | null) {
    return this.http.post<KanbanIssue[]>(`${this.apiUrl}/bind`, {issue_ids: issueIds, sprint_id: sprintId, assignee_id: assigneeId})
  }
}
