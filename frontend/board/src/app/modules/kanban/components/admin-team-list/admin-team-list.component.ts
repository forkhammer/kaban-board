import {Component, OnDestroy, OnInit} from '@angular/core';
import {TeamService} from "../../services/team.service";
import {finalize, Subject} from "rxjs";
import {takeUntil} from "rxjs/operators";
import {Team} from "../../models/team";
import { faPen, faTrash } from '@fortawesome/free-solid-svg-icons'

@Component({
  selector: 'app-admin-team-list',
  templateUrl: './admin-team-list.component.html',
  styleUrls: ['./admin-team-list.component.scss']
})
export class AdminTeamListComponent implements OnInit, OnDestroy {
  faPen = faPen
  faTrash = faTrash

  private destroy$ = new Subject()
  public teams: Team[] = []
  public isLoading = true

  constructor(
    private teamService: TeamService
  ) {
  }

  ngOnInit() {
    this.teamService.all().pipe(
      finalize(() => this.isLoading = false),
      takeUntil(this.destroy$)
    ).subscribe(data => {
      this.teams = data as Team[]
    })
  }

  ngOnDestroy() {
    this.destroy$.next(null)
    this.destroy$.complete()
  }

  trackByTeams(_: number, team: Team): number {
    return team.id
  }

  addTeam() {
    const existNew = this.teams.find(t => !Boolean(t.id))
    if (!existNew) {
      this.teams.push({
        id: 0,
        title: '',
        groups: []
      })
    }
  }

  catchOnDelete(team: Team) {
    this.teams.splice(this.teams.indexOf(team), 1)
  }
}
