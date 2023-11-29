import { Pipe, PipeTransform } from '@angular/core';
import { KanbanIssue } from '../models/kanban-issue';

@Pipe({
  name: 'filterIssuesByText'
})
export class FilterIssuesByTextPipe implements PipeTransform {

  transform(issues: KanbanIssue[], text: string | null): KanbanIssue[] {
    const lowerText = text ? text.toLowerCase() : null
    const issueId = text ? text.trim().replace('#', '') : null

    return issues.filter(issue => {
      return lowerText
          ? (issue.title.toLowerCase().indexOf(lowerText) > -1) || (issue.iid === issueId)
          : true
    })
  }

}
