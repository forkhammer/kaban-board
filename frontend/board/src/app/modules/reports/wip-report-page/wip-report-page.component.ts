import { Component, DestroyRef, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup } from '@angular/forms';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { BehaviorSubject, combineLatestWith, distinctUntilChanged, filter, map, switchMap } from 'rxjs';
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

interface IntervalOption {
  value: string;
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
  reload$ = new BehaviorSubject<null>(null);

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

    const interval$ = this.route.queryParams.pipe(map(p => p['interval']), distinctUntilChanged());
    const startDate$ = this.route.queryParams.pipe(map(p => p['start_date']), distinctUntilChanged());
    const endDate$ = this.route.queryParams.pipe(map(p => p['end_date']), distinctUntilChanged());
    const team$ = this.route.queryParams.pipe(map(p => p['team']), distinctUntilChanged());
    const user$ = this.route.queryParams.pipe(map(p => p['user']), distinctUntilChanged());

    // Sync URL -> form
    interval$.pipe(
      combineLatestWith(startDate$, endDate$, team$, user$),
      takeUntilDestroyed()
    ).subscribe(([interval, startDate, endDate, team, user]) => {
      this.form.patchValue({
        interval: interval ?? this.form.value.interval,
        start_date: startDate ?? this.form.value.start_date,
        end_date: endDate ?? this.form.value.end_date,
        team: team ?? null,
        user: user ?? null,
      }, { emitEvent: false });
    });

    // Sync form -> URL
    this.form.valueChanges.pipe(
      distinctUntilChanged(isEqual),
      takeUntilDestroyed()
    ).subscribe(data => {
      this.router.navigate([], { queryParams: data, queryParamsHandling: 'merge' });
    });

    // Load data when any param changes
    interval$.pipe(
      combineLatestWith(startDate$, endDate$, team$, user$, this.reload$),
      map(([interval, startDate, endDate, team, user]) => ({ interval, startDate, endDate, team, user })),
      distinctUntilChanged(isEqual),
      filter(p => !!(p.interval ?? this.form.value.interval) && !!(p.startDate ?? this.form.value.start_date) && !!(p.endDate ?? this.form.value.end_date)),
      switchMap(p => {
        this.loading = true;
        const v = this.form.value;
        return this.reportService.getWipReport({
          interval: p.interval ?? v.interval,
          start_date: p.startDate ?? v.start_date,
          end_date: p.endDate ?? v.end_date,
          team_id: (p.team ?? v.team) ? Number(p.team ?? v.team) : undefined,
          user_id: (p.user ?? v.user) ? Number(p.user ?? v.user) : undefined,
        }).pipe(catchErrorMessages(this.toast));
      }),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(report => {
      this.loading = false;
      if (report) {
        this.report = report;
        this.rebuildChart(report);
      }
    });

    // When params cleared -> clear chart
    interval$.pipe(
      combineLatestWith(startDate$, endDate$),
      filter(([interval, startDate, endDate]) => !interval && !startDate && !endDate),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(() => {
      this.report = null;
      this.graph = { ...this.graph, data: [] };
    });

    // Theme changes -> update layout
    this.themeService.theme$.pipe(
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(theme => {
      this.graph.layout = this.buildLayout(theme);
    });
  }

  private rebuildChart(report: WipReport): void {
    const theme = this.themeService.theme$.value;

    this.graph = {
      ...this.graph,
      data: [
        {
          x: report.data_points.map(p => p.label),
          y: report.data_points.map(p => p.wip_count),
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
        showgrid: true, gridcolor: gridColor, tickfont: { color: textColor }, type: 'category'
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
    this.reload$.next(null);
  }
}
