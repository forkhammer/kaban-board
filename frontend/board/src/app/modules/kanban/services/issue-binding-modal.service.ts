import { Injectable, inject } from '@angular/core';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { IssueBindingModalComponent } from '../components/issue-binding-modal/issue-binding-modal.component';
import { KanbanIssue } from '../models/kanban-issue';
import { User } from '../models/user';

export interface IssueBindingModalData {
  assigneeId: number | null;
  sprintId: number| null;
}

@Injectable({
  providedIn: 'root'
})
export class IssueBindingModalService {
  private modal = inject(NgbModal);

  /**
   * Открывает модальное окно создания Issue Binding
   * @param data - содержит assignee и sprintId для предзаполнения
   * @returns Promise с созданным KanbanIssue или null при отмене
   */
  show(data: IssueBindingModalData): Promise<KanbanIssue | null> {
    const modalRef = this.modal.open(IssueBindingModalComponent, {
      container: 'app-root',
      centered: true,
      size: 'lg'
    });

    modalRef.componentInstance.init(data);
    return modalRef.result as Promise<KanbanIssue | null>;
  }
}
