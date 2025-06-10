import {Injectable, Injector} from '@angular/core';
import {BaseService} from "../../core/services/base.service";
import { KanbanIssue } from '../models/kanban-issue';

@Injectable({
  providedIn: 'root'
})
export class IssueService extends BaseService<KanbanIssue>{
  public override usePagination = false

  constructor(protected override injector: Injector) {
    super(injector)
    this.apiUrl = this.config.apiUrl + '/issue'
  }

  formatter(item: KanbanIssue) {
    return `#${item.projectName}${item.iid} - ${item.title}`
  }
}
