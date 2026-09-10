import {Component, Input} from '@angular/core';
import {KanbanUser} from "../../models/kanban-user";
import {UserWorkload} from "../../../reports/models/report";

@Component({
    selector: 'app-user-card',
    templateUrl: './user-card.component.html',
    styleUrls: ['./user-card.component.scss'],
    standalone: false
})
export class UserCardComponent {
  @Input() user!: KanbanUser
  @Input() selected: boolean = false
  @Input() workload?: UserWorkload;

  get loadPercent(): number | null {
    if (!this.workload) return null;
    return this.workload.capacity > 0 ? Math.round(this.workload.planned / this.workload.capacity * 100) : 0;
  }

  get loadColorClass(): 'success' | 'warning' | 'danger' | null {
    const p = this.loadPercent;
    if (p === null) return null;
    if (p <= 80) return 'success';
    if (p <= 100) return 'warning';
    return 'danger';
  }
}
