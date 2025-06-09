import { Component, Input } from '@angular/core';
import { KanbanIssue } from '../../models/kanban-issue';

@Component({
  selector: 'app-issue-table-row',
  standalone: false,
  templateUrl: './issue-table-row.component.html',
  styleUrl: './issue-table-row.component.scss'
})
export class IssueTableRowComponent {
  @Input() issue!: KanbanIssue
}
