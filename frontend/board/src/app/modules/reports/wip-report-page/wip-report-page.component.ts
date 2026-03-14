import { Component, DestroyRef, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup } from '@angular/forms';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { distinctUntilChanged, filter, map, switchMap } from 'rxjs';
import * as PlotlyJS from 'plotly.js-dist-min';
import { PlotlyModule } from 'angular-plotly.js';
import { faArrowsRotate } from '@fortawesome/free-solid-svg-icons';

import { UiModule } from '../../ui/ui.module';
import { KanbanModule } from '../../kanban/kanban.module';
import { TeamService } from '../../kanban/services/team.service';
import { UserService } from '../../kanban/services/user.service';
import { ToastService } from '../../core/services/toast.service';
import { catchErrorMessages } from '../../core/tools/catch-error';
import { ThemeServiceService } from '../../ui/services/theme-service.service';
import { ReportService } from '../services/report.service';
import { WipReport } from '../models/report';
import { TitleService } from '../../core/services/title.service';
import { FaIconComponent } from '@fortawesome/angular-fontawesome';
import { ActivatedRoute, Router } from '@angular/router';
import { isEqual } from 'lodash';

PlotlyModule.plotlyjs = PlotlyJS;

type Interval = 'day' | 'week' | '2weeks' | 'month';

interface IntervalOption {
  value: Interval;
  label: string;
}

@Component({
  selector: 'app-wip-report-page',
  standalone: true,
  imports: [CommonModule, PlotlyModule, ReactiveFormsModule, UiModule, KanbanModule, FaIconComponent],
  templateUrl: './wip-report-page.component.html',
  styleUrl: './wip-report-page.component.scss'
})
export class WipReportPageComponent {
  private fb = inject(FormBuilder);
  private destroyRef = inject(DestroyRef);
  private toast = inject(ToastService);
  private reportService = inject(ReportService);
  private themeService = inject(ThemeServiceService);
  private title = inject(TitleService);
  private route = inject(ActivatedRoute);
  private router = inject(Router);

  teamService = inject(TeamService);
  userService = inject(UserService);

  form: FormGroup;
  report: WipReport | null = null;
  loading = false;
  faArrowsRotate = faArrowsRotate;

  intervalOptions: IntervalOption[] = [
    { value: 'day', label: 'День' },
    { value: 'week', label: 'Неделя' },
    { value: '2weeks', label: '2 недели' },
    { value: 'month', label: 'Месяц' },
  ];

  public graph: { data: any[]; layout: any; config: any } = {
    data: [],
    layout: this.buildLayout('light'),
    config: { responsive: true, displayModeBar: true, displaylogo: false }
  };

  constructor() {
    this.title.setTitleAndDescription('WIP Chart');

    const today = new Date();
    const monthAgo = new Date(today);
    monthAgo.setMonth(monthAgo.getMonth() - 1);

    this.form = this.fb.group({
      interval: ['week'],
      start_date: [this.formatDate(monthAgo)],
      end_date: [this.formatDate(today)],
      team: [null],
      user: [null],
    });

    // Restore state from URL query params (once on init)
    const queryParams = this.route.snapshot.queryParams;
    const patch: Record<string, any> = {};
    if (queryParams['interval']) patch['interval'] = queryParams['interval'];
    if (queryParams['start_date']) patch['start_date'] = queryParams['start_date'];
    if (queryParams['end_date']) patch['end_date'] = queryParams['end_date'];
    if (queryParams['team']) patch['team'] = queryParams['team'];
    if (queryParams['user']) patch['user'] = queryParams['user'];
    if (Object.keys(patch).length > 0) {
      this.form.patchValue(patch, { emitEvent: false });
    }

    // Sync form -> query params
    this.form.valueChanges.pipe(
      distinctUntilChanged(isEqual),
      takeUntilDestroyed()
    ).subscribe(data => {
      this.router.navigate([], { queryParams: data, queryParamsHandling: 'merge' });
    });

    // When backend params change -> load data
    const backendParams$ = this.form.valueChanges.pipe(
      map(v => ({ start_date: v.start_date, end_date: v.end_date, team: v.team, user: v.user })),
      distinctUntilChanged(isEqual),
      filter(v => !!v.start_date && !!v.end_date)
    );

    backendParams$.pipe(
      switchMap(v => {
        this.loading = true;
        return this.reportService.getWipReport({
          start_date: v.start_date,
          end_date: v.end_date,
          team_id: v.team ? Number(v.team) : undefined,
          user_id: v.user ? Number(v.user) : undefined,
        }).pipe(catchErrorMessages(this.toast));
      }),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(report => {
      this.loading = false;
      if (report) {
        this.report = report;
        this.rebuildChart();
      }
    });

    // When interval changes and report is already loaded -> re-render
    this.form.get('interval')!.valueChanges.pipe(
      distinctUntilChanged(),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(() => {
      if (this.report) this.rebuildChart();
    });

    // Theme changes -> update layout
    this.themeService.theme$.pipe(
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(theme => {
      this.graph.layout = this.buildLayout(theme);
    });

    // Initial load
    const v = this.form.value;
    if (v.start_date && v.end_date) {
      this.loading = true;
      this.reportService.getWipReport({
        start_date: v.start_date,
        end_date: v.end_date,
        team_id: v.team ? Number(v.team) : undefined,
        user_id: v.user ? Number(v.user) : undefined,
      }).pipe(
        catchErrorMessages(this.toast),
        takeUntilDestroyed(this.destroyRef)
      ).subscribe(report => {
        this.loading = false;
        if (report) {
          this.report = report;
          this.rebuildChart();
        }
      });
    }
  }

  private rebuildChart(): void {
    if (!this.report) return;
    const interval: Interval = this.form.value.interval || 'week';
    const theme = this.themeService.theme$.value;
    const grouped = this.groupByInterval(this.report.data_points, interval);

    this.graph = {
      ...this.graph,
      data: [
        {
          x: grouped.map(p => p.date),
          y: grouped.map(p => p.wip_count),
          type: 'scatter',
          mode: 'lines+markers',
          name: 'WIP',
          line: { color: '#1151F3', width: 2 },
          marker: { size: 6 }
        }
      ],
      layout: this.buildLayout(theme)
    };
  }

  private groupByInterval(dataPoints: WipReport['data_points'], interval: Interval): { date: string; wip_count: number }[] {
    if (interval === 'day') {
      return dataPoints.map(p => ({ date: p.date, wip_count: p.wip_count }));
    }

    const buckets = new Map<string, number[]>();
    for (const point of dataPoints) {
      const key = this.getBucketKey(point.date, interval);
      if (!buckets.has(key)) buckets.set(key, []);
      buckets.get(key)!.push(point.wip_count);
    }

    return Array.from(buckets.entries())
      .sort(([a], [b]) => a.localeCompare(b))
      .map(([date, values]) => ({
        date,
        wip_count: Math.round(values.reduce((s, v) => s + v, 0) / values.length)
      }));
  }

  private getBucketKey(dateStr: string, interval: Interval): string {
    const d = new Date(dateStr + 'T00:00:00');
    if (interval === 'week') {
      return this.formatDate(this.getMonday(d));
    }
    if (interval === '2weeks') {
      const monday = this.getMonday(d);
      const epoch = new Date('2000-01-03'); // первый понедельник эпохи
      const weeksSinceEpoch = Math.floor((monday.getTime() - epoch.getTime()) / (7 * 24 * 60 * 60 * 1000));
      const biWeekStart = new Date(epoch.getTime() + Math.floor(weeksSinceEpoch / 2) * 14 * 24 * 60 * 60 * 1000);
      return this.formatDate(biWeekStart);
    }
    if (interval === 'month') {
      return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-01`;
    }
    return dateStr;
  }

  private getMonday(date: Date): Date {
    const d = new Date(date);
    const day = d.getDay();
    const diff = d.getDate() - day + (day === 0 ? -6 : 1);
    d.setDate(diff);
    return d;
  }

  private formatDate(date: Date): string {
    return date.toISOString().slice(0, 10);
  }

  private buildLayout(theme: string): any {
    const isDark = theme === 'dark';
    const textColor = isDark ? '#dee2e6' : '#333333';
    const gridColor = isDark ? '#495057' : '#e0e0e0';
    const plotBg = isDark ? '#2b3035' : '#fafafa';
    const paperBg = isDark ? '#212529' : '#ffffff';
    const legendBg = isDark ? 'rgba(33, 37, 41, 0.8)' : 'rgba(255, 255, 255, 0.8)';

    return {
      title: { text: 'WIP — задачи в работе', font: { size: 20, color: textColor } },
      xaxis: {
        title: { text: 'Дата', font: { color: textColor } },
        showgrid: true, gridcolor: gridColor, tickfont: { color: textColor }
      },
      yaxis: {
        title: { text: 'Количество задач', font: { color: textColor } },
        showgrid: true, gridcolor: gridColor, tickfont: { color: textColor },
        rangemode: 'tozero'
      },
      legend: {
        x: 0.02, y: 0.98, bgcolor: legendBg,
        bordercolor: gridColor, borderwidth: 1, font: { color: textColor }
      },
      plot_bgcolor: plotBg,
      paper_bgcolor: paperBg,
      margin: { l: 60, r: 40, t: 60, b: 60 }
    };
  }

  reload(): void {
    const v = this.form.value;
    if (!v.start_date || !v.end_date) return;
    this.loading = true;
    this.reportService.getWipReport({
      start_date: v.start_date,
      end_date: v.end_date,
      team_id: v.team ? Number(v.team) : undefined,
      user_id: v.user ? Number(v.user) : undefined,
    }).pipe(
      catchErrorMessages(this.toast),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(report => {
      this.loading = false;
      if (report) {
        this.report = report;
        this.rebuildChart();
      }
    });
  }
}
