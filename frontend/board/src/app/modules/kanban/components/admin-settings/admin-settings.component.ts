import { Component, inject } from '@angular/core';
import { KanbanSettingsService } from '../../services/kanban-settings.service';
import { FormBuilder, FormGroup } from '@angular/forms';
import { LabelService } from '../../services/label.service';
import { distinctUntilChanged, filter, switchMap } from 'rxjs';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Component({
    selector: 'app-admin-settings',
    templateUrl: './admin-settings.component.html',
    styleUrls: ['./admin-settings.component.scss'],
    standalone: false
})
export class AdminSettingsComponent {
  private settingsServie = inject(KanbanSettingsService)
  private fb = inject(FormBuilder)
  public labelService = inject(LabelService)

  form: FormGroup
  private isUpdate = false

  constructor() {
    this.form = this.fb.group({
      taskTypeLabels: [[]]
    })

    this.settingsServie.getKanbanSettings().pipe(
      takeUntilDestroyed()
    ).subscribe(data => {
      this.isUpdate = true
      this.form.patchValue(data)
      this.isUpdate = false
    })

    this.form.get('taskTypeLabels')?.valueChanges.pipe(
      filter(_ => !this.isUpdate),
      distinctUntilChanged(),
      switchMap(data => this.settingsServie.saveTaskTypeLabels(data as string[])),
      takeUntilDestroyed()
    ).subscribe(data => {})
  }

}
