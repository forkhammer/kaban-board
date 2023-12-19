import { HttpClient } from "@angular/common/http";
import { Inject, Injectable } from "@angular/core";
import { CoreConfig, CoreConfigService } from "../../core/config";
import { map, shareReplay } from "rxjs";
import { KanbanSettings } from "../models/settings";

@Injectable({providedIn: 'root'})
export class KanbanSettingsService {
  constructor(
    private http: HttpClient,
    @Inject(CoreConfigService) private config: CoreConfig,
  ) {}

  getSettings() {
    return this.http.get(`${this.config.apiUrl}/settings`).pipe(
      map(data => data as KanbanSettings),
      shareReplay()
    )
  }
}
