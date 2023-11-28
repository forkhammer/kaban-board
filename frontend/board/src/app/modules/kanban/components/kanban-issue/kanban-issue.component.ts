import {Component, Input} from '@angular/core';
import {KanbanIssue} from "../../models/kanban-issue";

@Component({
  selector: 'app-kanban-issue',
  templateUrl: './kanban-issue.component.html',
  styleUrls: ['./kanban-issue.component.scss']
})
export class KanbanIssueComponent {
  @Input() issue: KanbanIssue | null = null
}
