import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { Label } from '../../models/kanban-label';
import { FormBuilder, FormGroup } from '@angular/forms';
import { Subject } from 'rxjs';
import { LabelService } from '../../services/label.service';
import {faPen, faTrash, faFloppyDisk} from "@fortawesome/free-solid-svg-icons";
import { catchErrorMessages } from 'src/app/modules/core/tools/catch-error';
import { ToastService } from 'src/app/modules/core/services/toast.service';

@Component({
  selector: 'app-admin-label-card',
  templateUrl: './admin-label-card.component.html',
  styleUrls: ['./admin-label-card.component.scss']
})
export class AdminLabelCardComponent implements OnInit, OnDestroy {
  @Input() label!: Label
  form: FormGroup
  private destroy$ = new Subject()
  protected readonly faPen = faPen
  protected readonly faTrash = faTrash
  protected readonly faFloppyDisk = faFloppyDisk
  public isEdit = false

  constructor(
    private fb: FormBuilder,
    public labelService: LabelService,
    private toast: ToastService
  ) {
    this.form = this.fb.group({
      altName: [null]
    })
  }

  ngOnInit() {
    this.form.patchValue(this.label)
  }

  ngOnDestroy() {
    this.destroy$.next(null);
    this.destroy$.complete()
  }

  setEdit() {
    this.isEdit = true
  }

  save() {
    let item = Object.assign(this.label, this.form.value);
    this.labelService.save(item)
      .pipe(
        catchErrorMessages(this.toast),
      )
      .subscribe(data => {
        this.isEdit = false
      });
  }
}
