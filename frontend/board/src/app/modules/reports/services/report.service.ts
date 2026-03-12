import { inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { CoreConfigService } from '../../core/config';
import { BurndownReport } from '../models/report';

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
}
