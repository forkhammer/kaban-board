import {Component, OnDestroy} from '@angular/core';
import {Subject} from "rxjs";
import {FormBuilder, FormGroup, Validators} from "@angular/forms";
import {NgbActiveModal} from "@ng-bootstrap/ng-bootstrap";
import {KanbanColumn} from "../../models/kanban-column";
import {markControlIfDirty} from "../../../core/tools/forms";
import {KanbanColumnService} from "../../services/kanban-column.service";
import {ToastService} from "../../../core/services/toast.service";
import {catchErrorMessages} from "../../../core/tools/catch-error";
import {LabelService} from "../../services/label.service";
import {TeamService} from "../../services/team.service";

@Component({
  selector: 'app-kanban-column-modal',
  templateUrl: './kanban-column-modal.component.html',
  styleUrls: ['./kanban-column-modal.component.scss'],
})
export class KanbanColumnModalComponent implements OnDestroy {
  private destroy$ = new Subject();
  public isLoading = false;
  public form: FormGroup;
  public item: KanbanColumn | null = null;

  constructor(
    public modal: NgbActiveModal,
    private fb: FormBuilder,
    private columnService: KanbanColumnService,
    private toast: ToastService,
    public labelService: LabelService,
    public teamService: TeamService
  ) {
    this.form = this.fb.group({
      id: [null, Validators.required],
      name: ['', Validators.required],
      labels: [[], Validators.required],
      team_id: [null]
    })
  }

  ngOnDestroy() {
    this.destroy$.next(null);
    this.destroy$.complete();
  }

  init(column: KanbanColumn | null = null) {
    this.item = column
    if (column) {
      this.form.patchValue(column)
      markControlIfDirty(this.form)
    }
  }

  close() {
    this.modal.close(this.item);
  }

  save(e: MouseEvent) {
    let item: KanbanColumn = {
      id: 0,
      name: '',
      labels: [],
      team_id: null,
      order: 10
    };
    if (this.item) {
      item = Object.assign(item, this.item);
    }
    item = Object.assign(item, this.form.value);
    this.isLoading = true;
    this.columnService.save(item)
      .pipe(
        catchErrorMessages(this.toast, () => this.isLoading = false),
      )
      .subscribe(data => {
        this.modal.close(data);
        this.isLoading = false;
      });
    return false;
  }
}
