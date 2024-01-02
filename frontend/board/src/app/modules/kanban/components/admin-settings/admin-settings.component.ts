import { Component, OnDestroy, OnInit } from '@angular/core';
import { KanbanSettingsService } from '../../services/kanban-settings.service';
import { FormBuilder, FormGroup } from '@angular/forms';
import { LabelService } from '../../services/label.service';
import { Subject, distinctUntilChanged, filter, switchMap, takeUntil } from 'rxjs';

@Component({
  selector: 'app-admin-settings',
  templateUrl: './admin-settings.component.html',
  styleUrls: ['./admin-settings.component.scss']
})
export class AdminSettingsComponent implements OnInit, OnDestroy {

  form: FormGroup
  private destroy$ = new Subject()
  private isUpdate = false
  constructor(
    private settingsServie: KanbanSettingsService,
    private fb: FormBuilder,
    public labelService: LabelService
  ) {
    this.form = this.fb.group({
      taskTypeLabels: [[]]
    })
  }

  ngOnInit(): void {
    this.settingsServie.getKanbanSettings().pipe(
      takeUntil(this.destroy$)
    ).subscribe(data => {
      this.isUpdate = true
      this.form.patchValue(data)
      this.isUpdate = false
    })

    this.form.get('taskTypeLabels')?.valueChanges.pipe(
      filter(_ => !this.isUpdate),
      distinctUntilChanged(),
      switchMap(data => this.settingsServie.saveTaskTypeLabels(data as string[])),
      takeUntil(this.destroy$)
    ).subscribe(data => {})
  }

  ngOnDestroy(): void {
    this.destroy$.next(null)
    this.destroy$.complete()
  }
}
