import { Component, DestroyRef, EventEmitter, inject, Input, OnInit, Output } from '@angular/core';
import { BIND_STATUS_LABELS, BIND_STATUS_VALUES, ISSUE_PRIORITY_LABELS, ISSUE_PRIORITY_VALUES, KanbanIssue } from '../../models/kanban-issue';
import { FormBuilder, FormGroup } from '@angular/forms';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { BehaviorSubject, debounceTime, filter, switchMap } from 'rxjs';
import { isEqual } from 'lodash';
import { IssueService } from '../../services/issue.service';
import { catchErrorMessages } from 'src/app/modules/core/tools/catch-error';
import { ToastService } from 'src/app/modules/core/services/toast.service';
import { UserService } from '../../services/user.service';
import { ReleaseService } from '../../services/release.service';
import { EpicService } from '../../services/epic.service';
import { AccountService } from 'src/app/modules/core/services/account.service';
import {faUser} from '@fortawesome/free-regular-svg-icons';
import { Team } from '../../models/team';

@Component({
  selector: 'app-issue-table-row, [app-issue-table-row]',
  standalone: false,
  templateUrl: './issue-table-row.component.html',
  styleUrl: './issue-table-row.component.scss'
})
export class IssueTableRowComponent implements OnInit{
  fb = inject(FormBuilder)
  destroyRef = inject(DestroyRef)
  issueService = inject(IssueService)
  toast = inject(ToastService)
  userService = inject(UserService)
  releaseService = inject(ReleaseService)
  epicService = inject(EpicService)
  accountService = inject(AccountService)

  readonly BIND_STATUS_VALUES = BIND_STATUS_VALUES
  readonly BIND_STATUS_LABELS =  BIND_STATUS_LABELS
  readonly ISSUE_PRIORITY_VALUES = ISSUE_PRIORITY_VALUES
  readonly ISSUE_PRIORITY_LABELS = ISSUE_PRIORITY_LABELS
  readonly faUser = faUser

  private _issue!: KanbanIssue
  form: FormGroup
  @Output() unbind = new EventEmitter<number>()
  isAdmin = false
  assigneeFilter: Record<string, any> = {}
  team$ = new BehaviorSubject<Team | null | undefined>(null)

  @Input()
  set issue(value: KanbanIssue) {
    this._issue = value
    this.form.patchValue(this.getSaveData(value), {emitEvent: false})
  }

  get issue(): KanbanIssue {
    return this._issue
  }

  @Input() set team(value: Team | undefined | null) {
    this.team$.next(value)
  }

  constructor() {
    this.form = this.fb.group({
      estimateDev: [null],
      estimateQA: [null],
      bindStatus: [null],
      assignee: [null],
      comment: [null],
      priority: [null],
      release: [null],
      epic: [null],
    })

    this.accountService.isAdmin$.pipe(takeUntilDestroyed()).subscribe(data => this.isAdmin = data)

    this.team$.pipe(takeUntilDestroyed()).subscribe(data => {
      if (data) {
        this.assigneeFilter = {team_id: data.id}
      } else {
        this.assigneeFilter = {}
      }
    })
  }

  ngOnInit(): void {
    this.form.valueChanges.pipe(
      debounceTime(500),
      filter(data => {
        return !isEqual(data, this.getSaveData(this._issue))
      }),
      switchMap(data => {
        const query = Object.assign({
          id: this._issue.id,
          bindingId: this._issue.bindingId,
        }, this.getSaveData(this._issue), data)
        return this.issueService.save(query).pipe(
          catchErrorMessages(this.toast)
        )
      }),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(data => {
      Object.assign(this._issue, data)
    })
  }

  unbindIssue() {
    if (this._issue.bindingId && confirm('Удалить задачу из спринта?')) {
      this.issueService.unbindFromSprint(this._issue.id, this._issue.bindingId).pipe(
        catchErrorMessages(this.toast),
        takeUntilDestroyed(this.destroyRef)
      ).subscribe(_ => {
        this.unbind.emit(this._issue.bindingId!)
      })
    }
  }

  private getSaveData(issue: KanbanIssue) {
    return {
      estimateDev: issue.estimateDev,
      estimateQA: issue.estimateQA,
      bindStatus: issue.bindStatus,
      assignee: issue.assignee ? issue.assignee.id : null,
      comment: issue.comment,
      priority: issue.priority,
      release: issue.release ? issue.release.id : null,
      epic: issue.epic ? issue.epic.id : null,
    }
  }
}
