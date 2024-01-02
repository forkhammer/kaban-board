import { HttpClient } from "@angular/common/http";
import { Inject, Injectable } from "@angular/core";
import { CoreConfig, CoreConfigService } from "../../core/config";
import { map, shareReplay } from "rxjs";
import { ClientSettings, KanbanSettings } from "../models/settings";

@Injectable({providedIn: 'root'})
export class KanbanSettingsService {
  constructor(
    private http: HttpClient,
    @Inject(CoreConfigService) private config: CoreConfig,
  ) {}

  getClientSettings() {
    return this.http.get(`${this.config.apiUrl}/settings`).pipe(
      map(data => data as ClientSettings),
      shareReplay()
    )
  }

  getKanbanSettings() {
    return this.http.get(`${this.config.apiUrl}/kanban-settings`).pipe(
      map(data => data as KanbanSettings),
    )
  }

  saveTaskTypeLabels(labels: string[]) {
    return this.http.post(`${this.config.apiUrl}/kanban-settings/task-type-labels`, {labels})
  }
}
