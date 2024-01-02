import { Pipe, PipeTransform } from '@angular/core';
import {KanbanIssue} from "../models/kanban-issue";

@Pipe({
  name: 'filterUnusedByLabels'
})
export class FilterUnusedByLabelsPipe implements PipeTransform {

  transform(issues: KanbanIssue[] | undefined, labels: string[]): KanbanIssue[] {
    return issues ? issues.filter(issue => {
      return !issue.labels.find(label => {
        return labels.find(s => label.name.toLowerCase().includes(s.toLowerCase()))
      })
    }) : []
  }

}
