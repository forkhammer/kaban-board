import { Component, Input } from '@angular/core';
import { ISSUE_PRIORITY_LABELS, IssuePriority } from '../../models/kanban-issue';

@Component({
  selector: 'app-priority-badge',
  templateUrl: './priority-badge.component.html',
  styleUrl: './priority-badge.component.scss',
  standalone: false,
})
export class PriorityBadgeComponent {
  @Input() priority!: IssuePriority

  readonly ISSUE_PRIORITY_LABELS = ISSUE_PRIORITY_LABELS
}
