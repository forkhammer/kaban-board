import { Injectable } from "@angular/core";
import { BindIssueModalComponent } from "../components/bind-issue-modal/bind-issue-modal.component";
import { NgbModal } from "@ng-bootstrap/ng-bootstrap";
import { KanbanIssue } from "../models/kanban-issue";

@Injectable({
  providedIn: 'root'
})
export class BindIssueModalService {
  constructor(private modal: NgbModal) { }

    show() {
      return new Promise((resolve, reject) => {
        const ref = this.modal.open(BindIssueModalComponent, {container: 'app-root', centered: true, size: 'lg'});
        // ref.componentInstance.init(column);
        ref.result.then(
          (result: KanbanIssue) => {
            if (resolve) {
              resolve(result);
            }
          },
          () => {
            if (reject) {
              reject(null);
            }
          }
        );
      });
    }
}
