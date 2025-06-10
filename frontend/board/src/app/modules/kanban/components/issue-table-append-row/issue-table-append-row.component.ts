import { Component, EventEmitter, inject, Input, Output } from '@angular/core';
import { KanbanIssue } from '../../models/kanban-issue';
import { FormBuilder, FormGroup } from '@angular/forms';
import { IssueService } from '../../services/issue.service';
import { ProjectService } from '../../services/project.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { Sprint } from '../../models/sprint';
import { filter, switchMap } from 'rxjs';
import { catchErrorMessages } from 'src/app/modules/core/tools/catch-error';
import { ToastService } from 'src/app/modules/core/services/toast.service';

@Component({
  selector: 'app-issue-table-append-row',
  standalone: false,
  templateUrl: './issue-table-append-row.component.html',
  styleUrl: './issue-table-append-row.component.scss'
})
export class IssueTableAppendRowComponent {
  fb = inject(FormBuilder)
  issueService = inject(IssueService)
  projectService = inject(ProjectService)
  toast = inject(ToastService)

  @Input() index: number = 0
  @Input() sprint!: Sprint
  @Output() issue = new EventEmitter<[number, KanbanIssue]>()

  form: FormGroup
  issueFilter: Record<string, any> = {}

  constructor() {
    this.form = this.fb.group({
      projectId: [null],
      issueId: [null]
    })

    this.form.get('projectId')?.valueChanges.pipe(
      takeUntilDestroyed()
    ).subscribe((project) => {
      this.issueFilter = { project }
    })

    this.form.get('issueId')?.valueChanges.pipe(
      filter(Boolean),
      switchMap(issueId => {
        return this.issueService.bindToSprint(issueId, this.sprint.id).pipe(
          catchErrorMessages(this.toast)
        )
      }),
      takeUntilDestroyed()
    ).subscribe(issue => {
      this.issue.emit([this.index, issue])
    })
  }
}
