import {Component, Input, OnDestroy, OnInit} from '@angular/core';
import {Project} from "../../models/project";
import {FormBuilder, FormGroup} from "@angular/forms";
import {TeamService} from "../../services/team.service";
import {distinctUntilChanged, Subject, switchMap} from "rxjs";
import {takeUntil} from "rxjs/operators";
import {ProjectService} from "../../services/project.service";

@Component({
  selector: 'app-admin-project-card',
  templateUrl: './admin-project-card.component.html',
  styleUrls: ['./admin-project-card.component.scss']
})
export class AdminProjectCardComponent implements OnInit, OnDestroy {
  @Input() project!:Project
  form: FormGroup
  private destroy$ = new Subject()

  constructor(
    private fb: FormBuilder,
    public teamService: TeamService,
    public projectService: ProjectService
  ) {
    this.form = this.fb.group({
      team_id: [null]
    })
  }

  ngOnInit() {
    this.form.patchValue(this.project)

    this.form.get('team_id')?.valueChanges.pipe(
      distinctUntilChanged(),
      switchMap(data => this.projectService.setTeam(this.project.id, data)),
      takeUntil(this.destroy$)
    ).subscribe(project => {
      Object.assign(this.project, project)
    })
  }

  ngOnDestroy() {
    this.destroy$.next(null);
    this.destroy$.complete()
  }
}
