import {Component, DestroyRef, EventEmitter, inject, Input, OnInit, Output} from '@angular/core';
import {Team} from "../../models/team";
import {faPen, faTrash, faFloppyDisk} from "@fortawesome/free-solid-svg-icons";
import {FormBuilder, FormGroup, Validators} from "@angular/forms";
import {TeamService} from "../../services/team.service";
import {catchErrorMessages} from "../../../core/tools/catch-error";
import {ToastService} from "../../../core/services/toast.service";
import { GroupService } from '../../services/group.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Component({
    selector: 'app-admin-team-card',
    templateUrl: './admin-team-card.component.html',
    styleUrls: ['./admin-team-card.component.scss'],
    standalone: false
})
export class AdminTeamCardComponent implements OnInit {
  private fb = inject(FormBuilder)
  private teamService = inject(TeamService)
  private toast = inject(ToastService)
  public groupService = inject(GroupService)
  private destroyRef = inject(DestroyRef)

  @Input() team!: Team
  public isEdit = false
  public form: FormGroup
  @Output() public onDelete = new EventEmitter<Team>()
  protected readonly faPen = faPen
  protected readonly faTrash = faTrash
  protected readonly faFloppyDisk = faFloppyDisk

  constructor() {
    this.form = this.fb.group({
      title: ['', Validators.required],
      groups: [[]],
    })
  }

  ngOnInit() {
    this.form.patchValue(this.getFormData(this.team))
  }

  isNew() {
    return !Boolean(this.team.id)
  }

  delete() {
    if (this.team.id) {
      if (confirm('Удалить группу?')) {
        this.teamService.delete(this.team).pipe(
          catchErrorMessages(this.toast),
          takeUntilDestroyed(this.destroyRef),
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
      groups: []
    };
    if (this.team) {
      item = Object.assign(item, this.getFormData(this.team));
    }
    item = Object.assign(item, this.form.value);
    this.teamService.save(item)
      .pipe(
        catchErrorMessages(this.toast),
        takeUntilDestroyed(this.destroyRef)
      )
      .subscribe(data => {
        Object.assign(this.team, data)
        this.isEdit = false
        this.form.patchValue(this.getFormData(this.team))
      });
  }

  setEdit() {
    this.isEdit = true
  }

  getFormData(team: Team) {
    return {
      id: team.id,
      title: team.title,
      groups: team.groups.map(g => g.id),
    }
  }
}
