import { Pipe, PipeTransform } from '@angular/core';
import {KanbanUser} from "../models/kanban-user";

@Pipe({
  name: 'searchUsers'
})
export class SearchUsersPipe implements PipeTransform {

  transform(users: KanbanUser[], search: string): KanbanUser[] {
    return users.filter(user => user.name.toLowerCase().includes(search.toLowerCase()))
  }

}
