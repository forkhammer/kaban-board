import { Component, DestroyRef, inject, Input, OnInit } from '@angular/core';
import { Label } from '../../models/kanban-label';
import { FormBuilder, FormGroup } from '@angular/forms';
import { LabelService } from '../../services/label.service';
import {faPen, faTrash, faFloppyDisk} from "@fortawesome/free-solid-svg-icons";
import { catchErrorMessages } from 'src/app/modules/core/tools/catch-error';
import { ToastService } from 'src/app/modules/core/services/toast.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Component({
    selector: 'app-admin-label-card',
    templateUrl: './admin-label-card.component.html',
    styleUrls: ['./admin-label-card.component.scss'],
    standalone: false
})
export class AdminLabelCardComponent implements OnInit {
  private fb = inject(FormBuilder)
  public labelService = inject(LabelService)
  private toast = inject(ToastService)
  private destroyRef = inject(DestroyRef)

  @Input() label!: Label
  form: FormGroup
  protected readonly faPen = faPen
  protected readonly faTrash = faTrash
  protected readonly faFloppyDisk = faFloppyDisk
  public isEdit = false

  constructor() {
    this.form = this.fb.group({
      altName: [null],
      bindingStatus: [null],
      priority: [null]
    })
  }

  ngOnInit() {
    this.form.patchValue(this.label)
  }

  setEdit() {
    this.isEdit = true
  }

  save() {
    let item = Object.assign(this.label, this.form.value);
    this.labelService.save(item)
      .pipe(
        catchErrorMessages(this.toast),
        takeUntilDestroyed(this.destroyRef)
      )
      .subscribe(data => {
        this.isEdit = false
      });
  }
}
