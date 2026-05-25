import {Component, inject, Input} from '@angular/core';
import {KanbanIssue} from "../../models/kanban-issue";
import {faSquare} from "@fortawesome/free-regular-svg-icons";
import {faHammer, faBug, faEye} from "@fortawesome/free-solid-svg-icons";
import {IssueDetailModalService} from "../../services/issue-detail-modal.service";

@Component({
    selector: 'app-kanban-issue',
    templateUrl: './kanban-issue.component.html',
    styleUrls: ['./kanban-issue.component.scss'],
    standalone: false
})
export class KanbanIssueComponent {
  @Input() issue: KanbanIssue | null = null

  private issueDetailModal = inject(IssueDetailModalService)

  faSquare = faSquare
  faHammer = faHammer
  faBug = faBug
  faEye = faEye

  get assignee() {
    if (this.issue?.assignee) {
      return this.issue.assignee
    }
    return this.issue?.assignees.length ? this.issue?.assignees[0] : null
  }

  showDetail(event: MouseEvent): void {
    event.stopPropagation()
    if (this.issue) {
      this.issueDetailModal.show(this.issue)
    }
  }
}
