import { Pipe, PipeTransform } from '@angular/core';
import {KanbanUser} from "../models/kanban-user";

@Pipe({
  name: 'filterUsersByGroup'
})
export class FilterUsersByGroupPipe implements PipeTransform {

  transform(users: KanbanUser[], groupId: number | null): KanbanUser[] {
    return users.filter(user => {
      return groupId ? user.groups.find(g => g.id === groupId) !== undefined : user.groups.length === 0
    })
  }

}
