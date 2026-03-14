import { Component, DestroyRef, inject } from '@angular/core';
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { IssueBindingService } from '../../services/issue-binding.service';
import { SprintService } from '../../services/sprint.service';
import { Sprint } from '../../models/sprint';
import { catchError, of } from 'rxjs';
import { faPlay, faStop, faClock } from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'app-copy-to-sprint-modal',
  templateUrl: './copy-to-sprint-modal.component.html',
  styleUrl: './copy-to-sprint-modal.component.scss',
  standalone: false
})
export class CopyToSprintModalComponent {
  modal = inject(NgbActiveModal);
  issueBindingService = inject(IssueBindingService);
  sprintService = inject(SprintService);
  destroyRef = inject(DestroyRef);

  readonly faPlay = faPlay;
  readonly faStop = faStop;
  readonly faClock = faClock;

  sprints: Sprint[] = [];
  isLoading = false;

  private bindingId!: number;
  private currentSprintId!: number;

  init(data: { bindingId: number; teamId: number; currentSprintId: number }) {
    this.bindingId = data.bindingId;
    this.currentSprintId = data.currentSprintId;
    this.isLoading = true;
    this.sprintService.all({ team: data.teamId }).pipe(
      catchError(() => of([])),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(result => {
      this.sprints = (result as Sprint[]).filter(s => s.id !== this.currentSprintId);
      this.isLoading = false;
    });
  }

  selectSprint(sprint: Sprint) {
    this.issueBindingService.copy(this.bindingId, sprint.id).pipe(
      takeUntilDestroyed(this.destroyRef)
    ).subscribe({
      next: result => {
        this.modal.close(result);
      },
      error: () => {
      }
    });
  }

  close() {
    this.modal.dismiss();
  }
}
