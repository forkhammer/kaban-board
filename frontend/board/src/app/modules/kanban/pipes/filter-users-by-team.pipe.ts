import { Pipe, PipeTransform } from '@angular/core';
import {KanbanUser} from "../models/kanban-user";

@Pipe({
  name: 'filterUsersByTeam'
})
export class FilterUsersByTeamPipe implements PipeTransform {

  transform(users: KanbanUser[], teamId: number | null): KanbanUser[] {
    return users.filter(user => {
      return teamId ? user.teams.includes(teamId) : true
    })
  }

}
