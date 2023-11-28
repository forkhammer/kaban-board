import {Component, EventEmitter, Input, OnDestroy, OnInit, Output} from '@angular/core';
import {Team} from "../../models/team";
import {faPen, faTrash, faFloppyDisk} from "@fortawesome/free-solid-svg-icons";
import {FormBuilder, FormGroup, Validators} from "@angular/forms";
import {TeamService} from "../../services/team.service";
import {takeUntil} from "rxjs/operators";
import {Subject} from "rxjs";
import {KanbanColumn} from "../../models/kanban-column";
import {catchErrorMessages} from "../../../core/tools/catch-error";
import {ToastService} from "../../../core/services/toast.service";

@Component({
  selector: 'app-admin-team-card',
  templateUrl: './admin-team-card.component.html',
  styleUrls: ['./admin-team-card.component.scss']
})
export class AdminTeamCardComponent implements OnDestroy, OnInit {
  @Input() team!: Team
  public isEdit = false
  public form: FormGroup
  @Output() public onDelete = new EventEmitter<Team>()
  protected readonly faPen = faPen
  protected readonly faTrash = faTrash
  protected readonly faFloppyDisk = faFloppyDisk
  private destroy$ = new Subject()

  constructor(
    private fb: FormBuilder,
    private teamService: TeamService,
    private toast: ToastService
  ) {
    this.form = this.fb.group({
      title: ['', Validators.required]
    })
  }

  ngOnInit() {
    this.form.patchValue(this.team)
  }

  ngOnDestroy() {
    this.destroy$.next(null)
    this.destroy$.complete()
  }

  isNew() {
    return !Boolean(this.team.id)
  }

  delete() {
    if (this.team.id) {
      if (confirm('Удалить группу?')) {
        this.teamService.delete(this.team).pipe(
          takeUntil(this.destroy$)
        ).subscribe(_ => {
          this.onDelete.emit(this.team)
        })
      }
    } else {
      this.onDelete.emit(this.team)
    }
  }

  save() {
    let item: Team = {
      id: 0,
      title: '',
    };
    if (this.team) {
      item = Object.assign(item, this.team);
    }
    item = Object.assign(item, this.form.value);
    this.teamService.save(item)
      .pipe(
        catchErrorMessages(this.toast),
      )
      .subscribe(data => {
        Object.assign(this.team, data)
        this.isEdit = false
        this.form.patchValue(this.team)
      });
  }

  setEdit() {
    this.isEdit = true
  }
}
