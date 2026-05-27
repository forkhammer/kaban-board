import {Injectable, inject} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {CoreConfigService} from "../../core/config";
import {SprintUserSettings, SaveSprintUserSettingsRequest} from '../models/sprint-user-settings';

@Injectable({
  providedIn: 'root'
})
export class SprintUserSettingsService {
  private http = inject(HttpClient);
  private config = inject(CoreConfigService);
  private apiUrl = this.config.apiUrl + '/sprint';

  getSettings(sprintId: number) {
    return this.http.get<SprintUserSettings[]>(`${this.apiUrl}/${sprintId}/user-settings`);
  }

  saveSettings(sprintId: number, data: SaveSprintUserSettingsRequest) {
    return this.http.post<SprintUserSettings>(`${this.apiUrl}/${sprintId}/user-settings`, data);
  }

  deleteSettings(sprintId: number, userId: number) {
    return this.http.delete<void>(`${this.apiUrl}/${sprintId}/user-settings/${userId}`);
  }
}
