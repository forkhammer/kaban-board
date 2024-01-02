import { Component, OnDestroy, OnInit } from '@angular/core';
import { KanbanSettingsService } from '../../services/kanban-settings.service';
import { FormBuilder, FormGroup } from '@angular/forms';
import { LabelService } from '../../services/label.service';
import { Subject, distinctUntilChanged, switchMap, takeUntil } from 'rxjs';

@Component({
  selector: 'app-admin-settings',
  templateUrl: './admin-settings.component.html',
  styleUrls: ['./admin-settings.component.scss']
})
export class AdminSettingsComponent implements OnInit, OnDestroy {

  form: FormGroup
  private destroy$ = new Subject()
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
      this.form.patchValue(data)
    })

    this.form.get('taskTypeLabels')?.valueChanges.pipe(
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
