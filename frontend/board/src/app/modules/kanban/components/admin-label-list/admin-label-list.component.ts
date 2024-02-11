import { Component, OnDestroy, OnInit } from '@angular/core';
import { Subject, finalize, takeUntil } from 'rxjs';
import { Label } from '../../models/kanban-label';
import { LabelService } from '../../services/label.service';
import { FormBuilder, FormGroup } from '@angular/forms';

@Component({
  selector: 'app-admin-label-list',
  templateUrl: './admin-label-list.component.html',
  styleUrls: ['./admin-label-list.component.scss']
})
export class AdminLabelListComponent implements OnInit, OnDestroy {
  private destroy$ = new Subject()
  public labels: Label[] = []
  public isLoading = true
  public filterForm: FormGroup

  constructor(
    private labelService: LabelService,
    private fb: FormBuilder
  ) {
    this.filterForm = this.fb.group({
      search: ['']
    })
  }

  ngOnInit() {
    this.labelService.all().pipe(
      finalize(() => this.isLoading = false),
      takeUntil(this.destroy$)
    ).subscribe(data => {
      this.labels = data as Label[]
    })
  }

  ngOnDestroy() {
    this.destroy$.next(null)
    this.destroy$.complete()
  }

  trackByLabel(_: number, label: Label): string {
    return label.id
  }
}
