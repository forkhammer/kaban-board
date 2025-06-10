import { Component, EventEmitter, inject, Output } from '@angular/core';
import { KanbanIssue } from '../../models/kanban-issue';
import { FormBuilder, FormGroup } from '@angular/forms';
import { IssueService } from '../../services/issue.service';
import { ProjectService } from '../../services/project.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

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

  @Output() issue = new EventEmitter<KanbanIssue>()

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
  }
}
