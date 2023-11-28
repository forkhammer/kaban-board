import { Pipe, PipeTransform } from '@angular/core';
import {KanbanColumn} from "../models/kanban-column";

@Pipe({
  name: 'filterColumnsByTeam',
  pure: false
})
export class FilterColumnsByTeamPipe implements PipeTransform {

  transform(columns: KanbanColumn[], teamId: number | null): KanbanColumn[] {
    return columns.filter(column => {
      if (teamId) {
        return column.team_id === teamId
      } else {
        return column.team_id === null
      }
    })
  }

}
