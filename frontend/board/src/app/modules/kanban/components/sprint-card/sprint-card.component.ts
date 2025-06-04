import { Component, DestroyRef, EventEmitter, inject, Input, Output } from '@angular/core';
import { Sprint, SprintStatus } from '../../models/sprint';
import {faClock} from '@fortawesome/free-regular-svg-icons';
import {faPlay, faStop} from '@fortawesome/free-solid-svg-icons';
import { SprintModalServiceService as SprintModalService } from '../../services/sprint-modal.service';
import { SprintService } from '../../services/sprint.service';
import { ToastService } from 'src/app/modules/core/services/toast.service';
import { catchErrorMessages } from 'src/app/modules/core/tools/catch-error';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Component({
  selector: 'app-sprint-card',
  standalone: false,
  templateUrl: './sprint-card.component.html',
  styleUrl: './sprint-card.component.scss'
})
export class SprintCardComponent {
  sprintModal = inject(SprintModalService)
  sprintService = inject(SprintService)
  destroyRef = inject(DestroyRef)
  toast = inject(ToastService)

  @Input() sprint!: Sprint
  @Output() onDelete = new EventEmitter<Sprint>()

  SprintStatus = SprintStatus
  faClock = faClock
  faPlay = faPlay
  faStop = faStop

  change() {
    this.sprintModal.show(this.sprint.id)
      .then(sprint => {
        if (sprint) {
          Object.assign(this.sprint, sprint)
        }
      }).catch(() => {})
  }

  delete() {
    if (confirm('Удалить спринт?')) {
      this.sprintService.delete(this.sprint).pipe(
        catchErrorMessages(this.toast),
        takeUntilDestroyed(this.destroyRef),
      ).subscribe(_ => {
        this.onDelete.emit(this.sprint)
      })
    }
  }

  run() {
    if (confirm('Запустить спринт?')) {
      this.sprintService.run(this.sprint.id).pipe(
        catchErrorMessages(this.toast),
        takeUntilDestroyed(this.destroyRef),
      ).subscribe(data => {
        Object.assign(this.sprint, data)
      })
    }
  }

  complete() {
    if (confirm('Завершить спринт?')) {
      this.sprintService.complete(this.sprint.id).pipe(
        catchErrorMessages(this.toast),
        takeUntilDestroyed(this.destroyRef),
      ).subscribe(data => {
        Object.assign(this.sprint, data)
      })
    }
  }
}
