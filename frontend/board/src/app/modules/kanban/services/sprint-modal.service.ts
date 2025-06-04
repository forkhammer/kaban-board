import { inject, Injectable } from '@angular/core';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { Sprint } from '../models/sprint';
import { SprintModalComponent } from '../components/sprint-modal/sprint-modal.component';

@Injectable({
  providedIn: 'root'
})
export class SprintModalServiceService {
  modal = inject(NgbModal)

  constructor() { }

  show(sprintId?: number) {
    const ref = this.modal.open(SprintModalComponent, {centered: true});
    ref.componentInstance.init(sprintId)
    return ref.result as Promise<Sprint | null>
  }
}
