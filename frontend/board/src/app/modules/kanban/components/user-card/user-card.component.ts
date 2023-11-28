import {Component, Input} from '@angular/core';
import {KanbanUser} from "../../models/kanban-user";
import {KanbanColumn} from "../../models/kanban-column";

@Component({
  selector: 'app-user-card',
  templateUrl: './user-card.component.html',
  styleUrls: ['./user-card.component.scss']
})
export class UserCardComponent {
  @Input() user!: KanbanUser
  @Input() selected: boolean = false
  @Input() columns: KanbanColumn[] = []

  getIssuesCount() {
    const labels = this.columns.reduce((prev, column) => {
      prev.push(...column.labels)
      return prev
    }, [] as string[])

    const issues = this.user.issues.filter(issue => {
      return issue.labels.nodes.find(label => {
        return labels.find(s => label.title.toLowerCase().includes(s.toLowerCase()))
      })
    })

    return issues.length
  }
}
