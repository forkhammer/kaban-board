import {Injectable, Injector} from '@angular/core';
import {BaseService} from "../../core/services/base.service";
import {Project} from "../models/project";
import {map} from "rxjs/operators";

@Injectable({
  providedIn: 'root'
})
export class ProjectService extends BaseService<Project>{
  public override usePagination = false

  constructor(protected override injector: Injector) {
    super(injector)
    this.apiUrl = this.config.apiUrl + '/projects'
  }

  setTeam(projectId: number, teamId: number) {
    return this.http.post(`${this.apiUrl}/${projectId}/set_team`, {team_id: teamId}).pipe(
      map(data => data as Project)
    )
  }
}
