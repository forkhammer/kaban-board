import { Component, inject } from '@angular/core';
import { finalize  } from 'rxjs';
import { Label } from '../../models/kanban-label';
import { LabelService } from '../../services/label.service';
import { FormBuilder, FormGroup } from '@angular/forms';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Component({
    selector: 'app-admin-label-list',
    templateUrl: './admin-label-list.component.html',
    styleUrls: ['./admin-label-list.component.scss'],
    standalone: false
})
export class AdminLabelListComponent {
  private labelService = inject(LabelService)
  private fb = inject(FormBuilder)

  public labels: Label[] = []
  public isLoading = true
  public filterForm: FormGroup

  constructor() {
    this.filterForm = this.fb.group({
      search: ['']
    })
    this.labelService.all().pipe(
      finalize(() => this.isLoading = false),
      takeUntilDestroyed()
    ).subscribe(data => {
      this.labels = data as Label[]
    })
  }

  trackByLabel(_: number, label: Label): string {
    return label.id
  }
}
