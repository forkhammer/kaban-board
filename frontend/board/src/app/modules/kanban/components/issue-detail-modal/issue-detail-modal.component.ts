import { Component, DestroyRef, inject } from "@angular/core";
import { NgbActiveModal } from "@ng-bootstrap/ng-bootstrap";
import {
  BIND_STATUS_LABELS,
  ISSUE_PRIORITY_LABELS,
  KanbanIssue,
} from "../../models/kanban-issue";
import { IssueBindingService } from "../../services/issue-binding.service";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { catchErrorMessages } from "src/app/modules/core/tools/catch-error";
import { ToastService } from "src/app/modules/core/services/toast.service";
import { Pagination } from "src/app/modules/core/models/base";

@Component({
  selector: "app-issue-detail-modal",
  templateUrl: "./issue-detail-modal.component.html",
  styleUrl: "./issue-detail-modal.component.scss",
  standalone: false,
})
export class IssueDetailModalComponent {
  modal = inject(NgbActiveModal);
  issueBindingService = inject(IssueBindingService);
  destroyRef = inject(DestroyRef);
  toast = inject(ToastService);

  issue: KanbanIssue | null = null;
  bindings: KanbanIssue[] = [];

  readonly BIND_STATUS_LABELS = BIND_STATUS_LABELS;
  readonly ISSUE_PRIORITY_LABELS = ISSUE_PRIORITY_LABELS;

  init(issue: KanbanIssue): void {
    this.issue = issue;

    this.issueBindingService
      .list({ issue: this.issue.id })
      .pipe(catchErrorMessages(this.toast), takeUntilDestroyed(this.destroyRef))
      .subscribe((data) => {
        this.bindings = (data as Pagination<KanbanIssue>).results;
      });
  }

  close(): void {
    this.modal.dismiss();
  }
}
