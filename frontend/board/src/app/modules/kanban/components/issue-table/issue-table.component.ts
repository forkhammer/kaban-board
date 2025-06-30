import { Component, inject, Input } from '@angular/core';
import { BehaviorSubject, combineLatestWith, debounceTime, distinctUntilChanged, filter, switchMap, timer } from 'rxjs';
import { KanbanUser } from '../../models/kanban-user';
import { Team } from '../../models/team';
import { Sprint } from '../../models/sprint';
import { environment } from 'src/environments/environment';
import { isEqual } from 'lodash';
import { catchErrorMessages } from 'src/app/modules/core/tools/catch-error';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { KanbanIssue } from '../../models/kanban-issue';
import { IssueService } from '../../services/issue.service';
import { ToastService } from 'src/app/modules/core/services/toast.service';
import { faPlus } from '@fortawesome/free-solid-svg-icons'

@Component({
  selector: 'app-issue-table',
  standalone: false,
  templateUrl: './issue-table.component.html',
  styleUrl: './issue-table.component.scss'
})
export class IssueTableComponent {
  private issueService = inject(IssueService)
  private toast = inject(ToastService)

  faPlus = faPlus

  public user$ = new BehaviorSubject<KanbanUser | null | undefined>(null)
  public team$ = new BehaviorSubject<Team | null | undefined>(null)
  public sprint$ = new BehaviorSubject<Sprint | null | undefined>(null)
  private timer$ = timer(0, environment.autoUpdateIssuesMin * 60 * 1000)
  public issues: KanbanIssue[] = []
  public appendedIssues: (number | null)[] = []

  @Input() set user(value : KanbanUser | undefined | null) {
    this.user$.next(value)
  }

  @Input() set team(value: Team | undefined | null) {
    this.team$.next(value)
  }

  @Input() set sprint(value: Sprint | undefined | null) {
    this.sprint$.next(value)
  }

  constructor() {
    this.timer$.pipe(
      combineLatestWith(this.user$, this.team$, this.sprint$),
      filter(([_, user, team, sprint]) => {
        return !!team
      }),
      distinctUntilChanged(isEqual),
      debounceTime(1),
      switchMap(([_, user, team, sprint]) => {
        const query: Record<string, any> = {
          'team': team!.id
        }
        if (user) {
          query['assignee'] = user.id
        }
        if (sprint) {
          query['sprint'] = sprint.id
        }
        return this.issueService.list(query).pipe(
          catchErrorMessages(this.toast)
        )
      }),
      takeUntilDestroyed()
    ).subscribe(data => {
      this.issues = data as KanbanIssue[]
    })
  }

  appendIssue() {
    this.appendedIssues.push(null)
  }

  bindIssue(event: [number, KanbanIssue]) {
    this.appendedIssues.splice(event[0], 1)
    this.issues.push(event[1])
  }

  unbindIssue(bindingId: number) {
    this.issues = this.issues.filter(issue => issue.bindingId !== bindingId)
  }
}
