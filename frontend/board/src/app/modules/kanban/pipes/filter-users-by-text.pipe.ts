import { Pipe, PipeTransform } from '@angular/core';
import { KanbanUser } from '../models/kanban-user';

@Pipe({
  name: 'filterUsersByText'
})
export class FilterUsersByTextPipe implements PipeTransform {

  transform(users: KanbanUser[], text: string | null): KanbanUser[] {
    const lowerText = text ? text.toLowerCase() : null
    const issueId = text ? text.trim().replace('#', '') : null

    return users.filter(user => {
      return user.issues.find(issue => {
        return lowerText
          ? (issue.title.toLowerCase().indexOf(lowerText) > -1) || (issue.iid === issueId)
          : true
      }) !== undefined
    });
  }

}
