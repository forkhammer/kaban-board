import {Component, inject} from '@angular/core';
import {TeamService} from "../../services/team.service";
import {finalize} from "rxjs";
import {Team} from "../../models/team";
import { faPen, faTrash } from '@fortawesome/free-solid-svg-icons'
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Component({
    selector: 'app-admin-team-list',
    templateUrl: './admin-team-list.component.html',
    styleUrls: ['./admin-team-list.component.scss'],
    standalone: false
})
export class AdminTeamListComponent {
  private teamService = inject(TeamService)

  faPen = faPen
  faTrash = faTrash

  public teams: Team[] = []
  public isLoading = true

  constructor() {
    this.teamService.all().pipe(
      finalize(() => this.isLoading = false),
      takeUntilDestroyed()
    ).subscribe(data => {
      this.teams = data as Team[]
    })
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
