import { Injectable } from '@angular/core';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { CopyToSprintModalComponent } from '../components/copy-to-sprint-modal/copy-to-sprint-modal.component';
import { KanbanIssue } from '../models/kanban-issue';

@Injectable({
  providedIn: 'root'
})
export class CopyToSprintModalService {
  constructor(private modal: NgbModal) {}

  show(bindingId: number, teamId: number, currentSprintId: number): Promise<KanbanIssue | null> {
    return new Promise((resolve) => {
      const ref = this.modal.open(CopyToSprintModalComponent, { container: 'app-root', centered: true });
      ref.componentInstance.init({ bindingId, teamId, currentSprintId });
      ref.result.then(
        (result: KanbanIssue) => resolve(result),
        () => resolve(null)
      );
    });
  }
}
