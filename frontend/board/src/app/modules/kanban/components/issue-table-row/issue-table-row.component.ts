import { AfterViewChecked, Component, DestroyRef, ElementRef, EventEmitter, HostBinding, inject, Input, OnInit, Output, ViewChild } from '@angular/core';
import { BIND_STATUS_LABELS, BIND_STATUS_VALUES, ISSUE_PRIORITY_LABELS, ISSUE_PRIORITY_VALUES, KanbanIssue } from '../../models/kanban-issue';
import { FormBuilder, FormGroup } from '@angular/forms';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { BehaviorSubject, debounceTime, filter } from 'rxjs';
import { isEqual } from 'lodash';
import { IssueService } from '../../services/issue.service';
import { UserService } from '../../services/user.service';
import { ReleaseService } from '../../services/release.service';
import {faUser} from '@fortawesome/free-regular-svg-icons';
import { Team } from '../../models/team';
import { Sprint } from '../../models/sprint';
import { CopyToSprintModalService } from '../../services/copy-to-sprint-modal.service';
import { MoveToSprintModalService } from '../../services/move-to-sprint-modal.service';
import { faEllipsisVertical, faArrowUpRightFromSquare, faGripVertical } from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'app-issue-table-row, [app-issue-table-row]',
  standalone: false,
  templateUrl: './issue-table-row.component.html',
  styleUrl: './issue-table-row.component.scss'
})
export class IssueTableRowComponent implements OnInit, AfterViewChecked {
  fb = inject(FormBuilder)
  destroyRef = inject(DestroyRef)
  issueService = inject(IssueService)
  copyToSprintModalService = inject(CopyToSprintModalService)
  moveToSprintModalService = inject(MoveToSprintModalService)
  userService = inject(UserService)
  releaseService = inject(ReleaseService)

  readonly BIND_STATUS_VALUES = BIND_STATUS_VALUES
  readonly BIND_STATUS_LABELS =  BIND_STATUS_LABELS
  readonly ISSUE_PRIORITY_VALUES = ISSUE_PRIORITY_VALUES
  readonly ISSUE_PRIORITY_LABELS = ISSUE_PRIORITY_LABELS
  readonly faUser = faUser
  readonly faEllipsisVertical = faEllipsisVertical
  readonly faArrowUpRightFromSquare = faArrowUpRightFromSquare
  readonly faGripVertical = faGripVertical

  @ViewChild('commentEl') commentEl?: ElementRef<HTMLElement>
  @ViewChild('titleEl') titleEl?: ElementRef<HTMLElement>

  private _issue!: KanbanIssue
  private _commentPending = false
  private _titlePending = false
  form: FormGroup
  @Output() unbind = new EventEmitter<number>()
  @Output() save = new EventEmitter<KanbanIssue>()
  @Output() deleteIssue = new EventEmitter<KanbanIssue>()
  assigneeFilter: Record<string, any> = {}
  team$ = new BehaviorSubject<Team | null | undefined>(null)
  sprint$ = new BehaviorSubject<Sprint | null | undefined>(null)

  @Input()
  set issue(value: KanbanIssue) {
    this._issue = value
    this.form.patchValue(this.getSaveData(value), {emitEvent: false})
    this._commentPending = true
    this._titlePending = true
  }

  get issue(): KanbanIssue {
    return this._issue
  }

  @HostBinding('class.unplanned')
  get isUnplanned(): boolean {
    return this._issue?.is_unplanned ?? false
  }

  @Input() set team(value: Team | undefined | null) {
    this.team$.next(value)
  }

  @Input() set sprint(value: Sprint | undefined | null) {
    this.sprint$.next(value)
  }

  constructor() {
    this.form = this.fb.group({
      title: [null],
      estimateDev: [null],
      estimateQA: [null],
      bindStatus: [null],
      assignee: [null],
      comment: [null],
      priority: [null],
      release: [null],
      epic: [null],
    })

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
      filter(data => !isEqual(data, this.getSaveData(this._issue))),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(data => {
      const query = Object.assign({
        id: this._issue.id,
        bindingId: this._issue.bindingId,
        version: this._issue.version,
      }, this.getSaveData(this._issue), data)
      this.save.emit(query)
    })
  }

  canUpdate(): boolean {
    return this.issue.can_update ?? false
  }

  canManage(): boolean {
    return this.issue.can_manage ?? false
  }

  ngAfterViewChecked() {
    if (this._commentPending && this.commentEl) {
      this.commentEl.nativeElement.innerText = this._issue.comment ?? ''
      this._commentPending = false
    }
    if (this._titlePending && this.titleEl) {
      this.titleEl.nativeElement.innerText = this._issue.title ?? ''
      this._titlePending = false
    }
  }

  onCommentInput(event: Event) {
    const text = (event.target as HTMLElement).innerText
    this.form.get('comment')?.setValue(text, {emitEvent: true})
  }

  onTitleInput(event: Event) {
    const text = (event.target as HTMLElement).innerText
    this.form.get('title')?.setValue(text, {emitEvent: true})
  }

  copyIssue() {
    const sprint = this.sprint$.value
    if (!this._issue.bindingId || !sprint) return
    this.copyToSprintModalService.show(this._issue.bindingId, sprint.team_id, sprint.id)
  }

  async moveIssue() {
    const sprint = this.sprint$.value
    if (!this._issue.bindingId || !sprint) return
    const result = await this.moveToSprintModalService.show(this._issue.bindingId, sprint.team_id, sprint.id)
    if (result) {
      this.unbind.emit(this._issue.bindingId)
    }
  }

  unbindIssue() {
    if (this._issue.bindingId && confirm('Удалить задачу из спринта?')) {
      this.deleteIssue.emit(this._issue)
    }
  }

  setPlanned(planned: boolean) {
    this.save.emit({
      ...this.getSaveData(this._issue),
      id: this._issue.id,
      bindingId: this._issue.bindingId,
      version: this._issue.version,
      planned: planned,
    } as any)
  }

  private getSaveData(issue: KanbanIssue) {
    return {
      title: issue.title,
      estimateDev: issue.estimateDev,
      estimateQA: issue.estimateQA,
      bindStatus: issue.bindStatus,
      assignee: issue.assignee ? issue.assignee.id : null,
      comment: issue.comment,
      priority: issue.priority,
      release: issue.release ? issue.release.id : null,
      epic: issue.epic ? issue.epic.id : null,
      planned: issue.planned,
    }
  }
}
