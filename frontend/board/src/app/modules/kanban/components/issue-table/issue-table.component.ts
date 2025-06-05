import { Component, Input } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { KanbanUser } from '../../models/kanban-user';
import { Team } from '../../models/team';

@Component({
  selector: 'app-issue-table',
  standalone: false,
  templateUrl: './issue-table.component.html',
  styleUrl: './issue-table.component.scss'
})
export class IssueTableComponent {
  public user$ = new BehaviorSubject<KanbanUser | null | undefined>(null)
  public team$ = new BehaviorSubject<Team | null | undefined>(null)

  @Input() set user(value : KanbanUser | undefined | null) {
    this.user$.next(value)
  }

  @Input() set team(value: Team | undefined | null) {
    this.team$.next(value)
  }
}
