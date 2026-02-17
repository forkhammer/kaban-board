import { Injectable } from "@angular/core";
import { CatalogPath } from "../types";

@Injectable({
  providedIn: 'root'
})
export class UrlService {
  getSprintUrl(sprintId: number, teamId: number | null): CatalogPath {
    return {
      url: '/',
      query: {
        team: teamId,
        view: 'list',
        sprint: sprintId
      }
    }
  }
}
