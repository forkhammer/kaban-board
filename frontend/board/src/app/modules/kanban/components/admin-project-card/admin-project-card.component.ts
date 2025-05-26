import {Component, DestroyRef, inject, Input, OnInit} from '@angular/core';
import {Project} from "../../models/project";
import {FormBuilder, FormGroup} from "@angular/forms";
import {TeamService} from "../../services/team.service";
import {distinctUntilChanged, switchMap} from "rxjs";
import {ProjectService} from "../../services/project.service";
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Component({
    selector: 'app-admin-project-card',
    templateUrl: './admin-project-card.component.html',
    styleUrls: ['./admin-project-card.component.scss'],
    standalone: false
})
export class AdminProjectCardComponent implements OnInit {
  private fb = inject(FormBuilder)
  public teamService = inject(TeamService)
  public projectService = inject(ProjectService)
  private destroyRef = inject(DestroyRef)

  @Input() project!:Project
  form: FormGroup

  constructor() {
    this.form = this.fb.group({
      team_id: [null]
    })
  }

  ngOnInit() {
    this.form.patchValue(this.project)

    this.form.get('team_id')?.valueChanges.pipe(
      distinctUntilChanged(),
      switchMap(data => this.projectService.setTeam(this.project.id, data)),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(project => {
      Object.assign(this.project, project)
    })
  }

}
