import { Injectable } from "@angular/core";
import { BindIssueModalComponent } from "../components/bind-issue-modal/bind-issue-modal.component";
import { NgbModal } from "@ng-bootstrap/ng-bootstrap";
import { KanbanIssue } from "../models/kanban-issue";

@Injectable({
  providedIn: 'root'
})
export class BindIssueModalService {
  lastProjectId: number | null = null;

  constructor(private modal: NgbModal) { }

  show() {
    return new Promise<KanbanIssue | KanbanIssue[]>((resolve, reject) => {
      const ref = this.modal.open(BindIssueModalComponent, {container: 'app-root', centered: true, size: 'lg'});
      ref.componentInstance.selectedProjectId = this.lastProjectId;
      ref.result.then(
        (result: KanbanIssue | KanbanIssue[]) => {
          resolve(result);
        },
        () => {
          reject(null);
        }
      );
    });
  }

  saveProjectId(projectId: number | null) {
    this.lastProjectId = projectId;
  }
}
