import { Injectable } from '@angular/core';
import {NgbModal} from "@ng-bootstrap/ng-bootstrap";
import {KanbanColumn} from "../models/kanban-column";
import {KanbanColumnModalComponent} from "../components/kanban-column-modal/kanban-column-modal.component";

@Injectable({
  providedIn: 'root'
})
export class KanbanColumnModalService {

  constructor(private modal: NgbModal) { }

  show(column: KanbanColumn | null) {
    return new Promise((resolve, reject) => {
      const ref = this.modal.open(KanbanColumnModalComponent, {container: 'app-root'});
      ref.componentInstance.init(column);
      ref.result.then(
        (result: KanbanColumn) => {
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
