import { Injectable, inject } from '@angular/core';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { IssueDetailModalComponent } from '../components/issue-detail-modal/issue-detail-modal.component';
import { KanbanIssue } from '../models/kanban-issue';

@Injectable({
  providedIn: 'root'
})
export class IssueDetailModalService {
  private modal = inject(NgbModal);

  show(issue: KanbanIssue): void {
    const modalRef = this.modal.open(IssueDetailModalComponent, {
      centered: true,
    });
    modalRef.componentInstance.init(issue);
  }
}
