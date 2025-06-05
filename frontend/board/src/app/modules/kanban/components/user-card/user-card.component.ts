import {Component, Input} from '@angular/core';
import {KanbanUser} from "../../models/kanban-user";

@Component({
    selector: 'app-user-card',
    templateUrl: './user-card.component.html',
    styleUrls: ['./user-card.component.scss'],
    standalone: false
})
export class UserCardComponent {
  @Input() user!: KanbanUser
  @Input() selected: boolean = false

  getIssuesCount() {
    return 0
  }
}
