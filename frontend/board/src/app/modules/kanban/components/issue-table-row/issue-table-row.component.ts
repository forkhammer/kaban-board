import { Component, DestroyRef, inject, Input, OnInit } from '@angular/core';
import { KanbanIssue } from '../../models/kanban-issue';
import { FormBuilder, FormGroup } from '@angular/forms';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { debounceTime, distinctUntilChanged, switchMap } from 'rxjs';
import { isEqual } from 'lodash';
import { IssueService } from '../../services/issue.service';
import { catchErrorMessages } from 'src/app/modules/core/tools/catch-error';
import { ToastService } from 'src/app/modules/core/services/toast.service';

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

  private _issue!: KanbanIssue
  form: FormGroup

  @Input()
  set issue(value: KanbanIssue) {
    this._issue = value
    this.form.patchValue({
      estimateDev: value.estimateDev,
      estimateQA: value.estimateQA,
    })
  }

  get issue(): KanbanIssue {
    return this._issue
  }

  constructor() {
    this.form = this.fb.group({
      estimateDev: [null],
      estimateQA: [null],
    })
  }

  ngOnInit(): void {
    this.form.valueChanges.pipe(
      distinctUntilChanged(isEqual),
      debounceTime(10),
      switchMap(data => {
        const query = Object.assign({}, this._issue, data)
        return this.issueService.save(query).pipe(
          catchErrorMessages(this.toast)
        )
      }),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(data => {
      Object.assign(this._issue, data)
    })
  }
}
