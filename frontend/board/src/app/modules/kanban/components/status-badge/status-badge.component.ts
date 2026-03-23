import { Component, Input } from '@angular/core';
import { BIND_STATUS_LABELS, BindStatus } from '../../models/kanban-issue';

@Component({
  selector: 'app-status-badge',
  templateUrl: './status-badge.component.html',
  styleUrl: './status-badge.component.scss',
  standalone: false,
})
export class StatusBadgeComponent {
  @Input() status!: BindStatus

  readonly BIND_STATUS_LABELS = BIND_STATUS_LABELS;
}
