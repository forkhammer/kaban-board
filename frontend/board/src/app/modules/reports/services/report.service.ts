import { inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { CoreConfigService } from '../../core/config';
import { BurndownReport, BurnupReport, WipReport } from '../models/report';

@Injectable({
  providedIn: 'root'
})
export class ReportService {
  private http = inject(HttpClient);
  private config = inject(CoreConfigService);

  getBurndownReport(sprintId: number): Observable<BurndownReport> {
    return this.http.get<BurndownReport>(`${this.config.apiUrl}/reports/burndown`, {
      params: { sprint_id: sprintId.toString() }
    });
  }

  getBurnupReport(sprintId: number): Observable<BurnupReport> {
    return this.http.get<BurnupReport>(`${this.config.apiUrl}/reports/burnup`, {
      params: { sprint_id: sprintId.toString() }
    });
  }

  getWipReport(params: { start_date: string; end_date: string; interval: string; team_id?: number; user_id?: number }): Observable<WipReport> {
    const httpParams: Record<string, string> = {
      start_date: params.start_date,
      end_date: params.end_date,
      interval: params.interval,
    };
    if (params.team_id != null) httpParams['team_id'] = params.team_id.toString();
    if (params.user_id != null) httpParams['user_id'] = params.user_id.toString();
    return this.http.get<WipReport>(`${this.config.apiUrl}/reports/wip`, { params: httpParams });
  }
}
