import { Component, DestroyRef, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup } from '@angular/forms';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { BehaviorSubject, combineLatestWith, distinctUntilChanged, filter, map, switchMap } from 'rxjs';
import * as PlotlyJS from 'plotly.js-dist-min';
import { PlotlyModule } from 'angular-plotly.js';
import { FaIconComponent } from "@fortawesome/angular-fontawesome";
import { faArrowsRotate } from '@fortawesome/free-solid-svg-icons'

import { UiModule } from '../../ui/ui.module';
import { KanbanModule } from '../../kanban/kanban.module';
import { TeamService } from '../../kanban/services/team.service';
import { ToastService } from '../../core/services/toast.service';
import { catchErrorMessages } from '../../core/tools/catch-error';
import { ThemeServiceService } from '../../ui/services/theme-service.service';
import { ReportService } from '../services/report.service';
import { BurnupReport } from '../models/report';
import { TitleService } from '../../core/services/title.service';
import { ActivatedRoute, Router } from '@angular/router';
import { isEqual } from 'lodash';

PlotlyModule.plotlyjs = PlotlyJS;

@Component({
  selector: 'app-burnup-report-page',
  standalone: true,
  imports: [CommonModule, PlotlyModule, ReactiveFormsModule, UiModule, KanbanModule, FaIconComponent],
  templateUrl: './burnup-report-page.component.html',
  styleUrl: './burnup-report-page.component.scss'
})
export class BurnupReportPageComponent {
  private fb = inject(FormBuilder);
  private destroyRef = inject(DestroyRef);
  private toast = inject(ToastService);
  private reportService = inject(ReportService);
  private themeService = inject(ThemeServiceService);
  private title = inject(TitleService)
  private route = inject(ActivatedRoute)
  private router = inject(Router)

  teamService = inject(TeamService);

  form: FormGroup;
  sprintFilter: Record<string, any> = {};
  reload$ = new BehaviorSubject<null>(null);

  report: BurnupReport | null = null;
  loading = false;
  faArrowsRotate = faArrowsRotate

  // Dev stats
  currentScopeDev = 0;
  completedDev = 0;
  completedPercentDev = 0;

  // QA stats
  currentScopeQA = 0;
  completedQA = 0;
  completedPercentQA = 0;

  public devGraph: { data: any[]; layout: any; config: any } = {
    data: [],
    layout: this.buildLayout('light', 'Разработка'),
    config: {
      responsive: true,
      displayModeBar: true,
      displaylogo: false
    }
  };

  public qaGraph: { data: any[]; layout: any; config: any } = {
    data: [],
    layout: this.buildLayout('light', 'QA'),
    config: {
      responsive: true,
      displayModeBar: true,
      displaylogo: false
    }
  };

  constructor() {
    this.title.setTitleAndDescription('Burnup Chart')

    const team$ = this.route.queryParams.pipe(map(params => params['team']), distinctUntilChanged())
    const sprint$ = this.route.queryParams.pipe(map(params => params['sprint']), distinctUntilChanged())

    this.form = this.fb.group({
      team: [null],
      sprint: [null]
    });

    sprint$.pipe(
      combineLatestWith(team$),
      takeUntilDestroyed()
    ).subscribe(([sprint, team]) => {
      console.log('load form data', {team, sprint})
      this.form.patchValue({team, sprint}, {emitEvent: false})
    })

    this.form.valueChanges.pipe(
      distinctUntilChanged(isEqual),
      takeUntilDestroyed()
    ).subscribe(data => {
      console.log('set form data', data)
      this.router.navigate([], {queryParams: data, queryParamsHandling: 'merge'})
    })

    // When team changes, update sprint filter and reset sprint selection
    this.form.get('team')!.valueChanges.pipe(
      distinctUntilChanged(),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(teamId => {
      this.sprintFilter = teamId ? { team: teamId } : {};
      this.form.get('sprint')!.setValue(null);
    });

    // When sprint changes, load burnup data
    sprint$.pipe(
      distinctUntilChanged(),
      combineLatestWith(this.reload$),
      map(([sprintId, _]) => sprintId),
      filter(sprintId => sprintId != null),
      switchMap(sprintId => {
        this.loading = true;
        return this.reportService.getBurnupReport(sprintId).pipe(
          catchErrorMessages(this.toast)
        );
      }),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(report => {
      this.loading = false;
      if (report) {
        this.report = report;
        this.buildCharts(report);
      }
    });

    // When sprint is cleared
    sprint$.pipe(
      distinctUntilChanged(),
      filter(sprintId => sprintId == null),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(() => {
      this.report = null;
      this.devGraph.data = [];
      this.qaGraph.data = [];
      this.resetStats();
    });

    // React to theme changes to update chart colors
    this.themeService.theme$.pipe(
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(theme => {
      this.devGraph.layout = this.buildLayout(theme, 'Разработка');
      this.qaGraph.layout = this.buildLayout(theme, 'QA');
    });
  }

  private resetStats(): void {
    this.currentScopeDev = 0;
    this.completedDev = 0;
    this.completedPercentDev = 0;
    this.currentScopeQA = 0;
    this.completedQA = 0;
    this.completedPercentQA = 0;
  }

  private buildCharts(report: BurnupReport): void {
    const theme = this.themeService.theme$.value;

    // Build Dev chart
    this.buildDevChart(report, theme);

    // Build QA chart
    this.buildQAChart(report, theme);
  }

  private buildDevChart(report: BurnupReport, theme: string): void {
    const dates = report.data_points.map(dp => dp.date);
    const scopeDev = report.data_points.map(dp => dp.scope_dev);
    const completedDev = report.data_points.map(dp => dp.completed_dev);

    // Calculate stats from the last data point
    const lastDataPoint = report.data_points.length > 0
      ? report.data_points[report.data_points.length - 1]
      : null;

    if (lastDataPoint) {
      this.currentScopeDev = lastDataPoint.scope_dev;
      this.completedDev = lastDataPoint.completed_dev;
      this.completedPercentDev = lastDataPoint.scope_dev > 0
        ? Math.round((lastDataPoint.completed_dev / lastDataPoint.scope_dev) * 100)
        : 0;
    }

    this.devGraph.data = [
      {
        x: dates,
        y: scopeDev,
        type: 'scatter',
        mode: 'lines+markers',
        name: 'Scope',
        line: {
          color: '#dc3545',
          width: 2
        },
        marker: {
          size: 6
        }
      },
      {
        x: dates,
        y: completedDev,
        type: 'scatter',
        mode: 'lines+markers',
        name: 'Выполнено',
        line: {
          color: '#28a745',
          width: 2
        },
        marker: {
          size: 6
        }
      }
    ];

    this.devGraph.layout = this.buildLayout(theme, 'Разработка');
  }

  private buildQAChart(report: BurnupReport, theme: string): void {
    const dates = report.data_points.map(dp => dp.date);
    const scopeQA = report.data_points.map(dp => dp.scope_qa);
    const completedQA = report.data_points.map(dp => dp.completed_qa);

    // Calculate stats from the last data point
    const lastDataPoint = report.data_points.length > 0
      ? report.data_points[report.data_points.length - 1]
      : null;

    if (lastDataPoint) {
      this.currentScopeQA = lastDataPoint.scope_qa;
      this.completedQA = lastDataPoint.completed_qa;
      this.completedPercentQA = lastDataPoint.scope_qa > 0
        ? Math.round((lastDataPoint.completed_qa / lastDataPoint.scope_qa) * 100)
        : 0;
    }

    this.qaGraph.data = [
      {
        x: dates,
        y: scopeQA,
        type: 'scatter',
        mode: 'lines+markers',
        name: 'Scope',
        line: {
          color: '#dc3545',
          width: 2
        },
        marker: {
          size: 6
        }
      },
      {
        x: dates,
        y: completedQA,
        type: 'scatter',
        mode: 'lines+markers',
        name: 'Выполнено',
        line: {
          color: '#28a745',
          width: 2
        },
        marker: {
          size: 6
        }
      }
    ];

    this.qaGraph.layout = this.buildLayout(theme, 'QA');
  }

  private buildLayout(theme: string, type: string): any {
    const isDark = theme === 'dark';
    const textColor = isDark ? '#dee2e6' : '#333333';
    const gridColor = isDark ? '#495057' : '#e0e0e0';
    const plotBg = isDark ? '#2b3035' : '#fafafa';
    const paperBg = isDark ? '#212529' : '#ffffff';
    const legendBg = isDark ? 'rgba(33, 37, 41, 0.8)' : 'rgba(255, 255, 255, 0.8)';

    const titleText = type === 'Разработка' ? 'Dev Burnup' : 'QA Burnup';
    const yAxisText = type === 'Разработка' ? 'Story Points (Разработка)' : 'Story Points (QA)';

    return {
      title: {
        text: this.report ? `${this.report.sprint_title} - ${titleText}` : titleText,
        font: {
          size: 20,
          color: textColor
        }
      },
      xaxis: {
        title: { text: 'Дни спринта', font: { color: textColor } },
        showgrid: true,
        gridcolor: gridColor,
        tickfont: { color: textColor }
      },
      yaxis: {
        title: { text: yAxisText, font: { color: textColor } },
        showgrid: true,
        gridcolor: gridColor,
        tickfont: { color: textColor },
        rangemode: 'tozero'
      },
      legend: {
        x: 0.02,
        y: 0.98,
        bgcolor: legendBg,
        bordercolor: gridColor,
        borderwidth: 1,
        font: { color: textColor }
      },
      plot_bgcolor: plotBg,
      paper_bgcolor: paperBg,
      margin: {
        l: 60,
        r: 40,
        t: 60,
        b: 60
      }
    };
  }

  reload() {
    this.reload$.next(null)
  }
}
