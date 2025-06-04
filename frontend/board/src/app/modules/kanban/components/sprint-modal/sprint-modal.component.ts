import { Component, DestroyRef, inject } from '@angular/core';
import { SprintService } from '../../services/sprint.service';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { ToastService } from 'src/app/modules/core/services/toast.service';
import { SaveSprintRequest, Sprint } from '../../models/sprint';
import { catchErrorMessages } from 'src/app/modules/core/tools/catch-error';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';
import { catchError, EMPTY } from 'rxjs';
import { CommonModule, formatDate } from '@angular/common';
import { UiModule } from 'src/app/modules/ui/ui.module';
import { TeamService } from '../../services/team.service';

@Component({
  selector: 'app-sprint-modal',
  imports: [
    CommonModule,
    ReactiveFormsModule,
    UiModule
  ],
  templateUrl: './sprint-modal.component.html',
  styleUrl: './sprint-modal.component.scss'
})
export class SprintModalComponent {
  sprintService = inject(SprintService)
  fb = inject(FormBuilder)
  toast = inject(ToastService)
  destroyRef = inject(DestroyRef)
  modal = inject(NgbActiveModal)
  teamService = inject(TeamService)

  form: FormGroup
  sprint: Sprint | null = null
  isLoading = false
  errors: any

  constructor() {
    this.form = this.fb.group({
      title: [''],
      start_date: [formatDate(this.getStartSprintDate(), 'yyyy-MM-dd', 'en'), Validators.required],
      end_date: [formatDate(this.getEndSprintDate(), 'yyyy-MM-dd', 'en'), Validators.required],
      hours_per_user: [60],
      team_id: [null, Validators.required],
    })
  }

  init(sprintId?: number) {
    if (sprintId) {
      this.sprintService.get(sprintId).pipe(
        catchErrorMessages(this.toast),
        takeUntilDestroyed(this.destroyRef)
      ).subscribe(sprint => {
        this.sprint = sprint
        this.form.patchValue(sprint)
      })
    }
  }

  close() {
    this.modal.close(null);
  }

  submit() {
    this.isLoading = true
    const data: SaveSprintRequest = {
      id: this.sprint?.id,
      title: this.form.value.title,
      start_date: this.form.value.start_date,
      end_date: this.form.value.end_date,
      hours_per_user: this.form.value.hours_per_user,
      team_id: this.form.value.team_id
    }
    this.sprintService.save(data).pipe(
      catchError(err => {
        this.errors = err.error
        this.isLoading = false
        return EMPTY
      }),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(sprint => {
      this.isLoading = false
      this.modal.close(sprint)
    })
  }

  getStartSprintDate() {
    return new Date()
  }

  getEndSprintDate() {
    const dt = this.getStartSprintDate()
    dt.setDate(dt.getDate() + 14)
    return dt
  }
}
