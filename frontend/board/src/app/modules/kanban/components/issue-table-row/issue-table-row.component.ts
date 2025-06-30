import { Component, DestroyRef, EventEmitter, inject, Input, OnInit, Output } from '@angular/core';
import { BIND_STATUS_LABELS, BIND_STATUS_VALUES, KanbanIssue } from '../../models/kanban-issue';
import { FormBuilder, FormGroup } from '@angular/forms';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { debounceTime, distinctUntilChanged, switchMap } from 'rxjs';
import { isEqual } from 'lodash';
import { IssueService } from '../../services/issue.service';
import { catchErrorMessages } from 'src/app/modules/core/tools/catch-error';
import { ToastService } from 'src/app/modules/core/services/toast.service';
import { UserService } from '../../services/user.service';

@Component({
  selector: 'app-issue-table-row',
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

  readonly BIND_STATUS_VALUES = BIND_STATUS_VALUES

  private _issue!: KanbanIssue
  form: FormGroup
  @Output() unbind = new EventEmitter<number>()

  @Input()
  set issue(value: KanbanIssue) {
    this._issue = value
    this.form.patchValue({
      estimateDev: value.estimateDev,
      estimateQA: value.estimateQA,
      bindStatus: value.bindStatus,
      assignee: value.assignee ? value.assignee.id : null,
      comment: value.comment,
    })
  }

  get issue(): KanbanIssue {
    return this._issue
  }

  constructor() {
    this.form = this.fb.group({
      estimateDev: [null],
      estimateQA: [null],
      bindStatus: [null],
      assignee: [null],
      comment: [null]
    })
  }

  ngOnInit(): void {
    this.form.valueChanges.pipe(
      distinctUntilChanged(isEqual),
      debounceTime(500),
      switchMap(data => {
        const query = Object.assign({}, this._issue, data)
        console.log(query)
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
}
