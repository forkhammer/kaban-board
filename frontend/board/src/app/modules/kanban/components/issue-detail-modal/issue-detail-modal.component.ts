import { Component, inject } from '@angular/core';
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';
import {
  BIND_STATUS_LABELS,
  ISSUE_PRIORITY_LABELS,
  KanbanIssue
} from '../../models/kanban-issue';

@Component({
  selector: 'app-issue-detail-modal',
  templateUrl: './issue-detail-modal.component.html',
  styleUrl: './issue-detail-modal.component.scss',
  standalone: false
})
export class IssueDetailModalComponent {
  modal = inject(NgbActiveModal);

  issue: KanbanIssue | null = null;

  readonly BIND_STATUS_LABELS = BIND_STATUS_LABELS;
  readonly ISSUE_PRIORITY_LABELS = ISSUE_PRIORITY_LABELS;

  init(issue: KanbanIssue): void {
    this.issue = issue;
  }

  close(): void {
    this.modal.dismiss();
  }
}
