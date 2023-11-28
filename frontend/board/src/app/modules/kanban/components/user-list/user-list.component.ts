import {Component, Input} from '@angular/core';
import {KanbanUser} from "../../models/kanban-user";

@Component({
  selector: 'app-user-list',
  templateUrl: './user-list.component.html',
  styleUrls: ['./user-list.component.scss']
})
export class UserListComponent {
  @Input() users: KanbanUser[] = []

  trackByUser(index: number, user: KanbanUser) {
    return user.id
  }
}
