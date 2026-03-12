import { Component, DestroyRef, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup } from '@angular/forms';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { BehaviorSubject, combineLatestWith, distinctUntilChanged, filter, map, switchMap } from 'rxjs';
import * as PlotlyJS from 'plotly.js-dist-min';
import { PlotlyModule } from 'angular-plotly.js';
import { faArrowsRotate } from '@fortawesome/free-solid-svg-icons'

import { UiModule } from '../../ui/ui.module';
import { KanbanModule } from '../../kanban/kanban.module';
import { TeamService } from '../../kanban/services/team.service';
import { ToastService } from '../../core/services/toast.service';
import { catchErrorMessages } from '../../core/tools/catch-error';
import { ThemeServiceService } from '../../ui/services/theme-service.service';
import { ReportService } from '../services/report.service';
import { BurndownReport } from '../models/report';
import { TitleService } from '../../core/services/title.service';
import { FaIconComponent } from "@fortawesome/angular-fontawesome";

PlotlyModule.plotlyjs = PlotlyJS;

@Component({
  selector: 'app-burndown-report-page',
  standalone: true,
  imports: [CommonModule, PlotlyModule, ReactiveFormsModule, UiModule, KanbanModule, FaIconComponent],
  templateUrl: './burndown-report-page.component.html',
  styleUrl: './burndown-report-page.component.scss'
})
export class BurndownReportPageComponent {
  private fb = inject(FormBuilder);
  private destroyRef = inject(DestroyRef);
  private toast = inject(ToastService);
  private reportService = inject(ReportService);
  private themeService = inject(ThemeServiceService);
  private title = inject(TitleService)

  teamService = inject(TeamService);

  form: FormGroup;
  sprintFilter: Record<string, any> = {};
  reload$ = new BehaviorSubject<null>(null);

  report: BurndownReport | null = null;
  loading = false;
  faArrowsRotate = faArrowsRotate

  // Dev stats
  totalScopeDev = 0;
  remainingDev = 0;
  completedPercentDev = 0;

  // QA stats
  totalScopeQA = 0;
  remainingQA = 0;
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
    this.title.setTitleAndDescription('Burndown Chart')

    this.form = this.fb.group({
      team: [null],
      sprint: [null]
    });

    // When team changes, update sprint filter and reset sprint selection
    this.form.get('team')!.valueChanges.pipe(
      distinctUntilChanged(),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(teamId => {
      this.sprintFilter = teamId ? { team: teamId } : {};
      this.form.get('sprint')!.setValue(null);
    });

    // When sprint changes, load burndown data
    this.form.get('sprint')!.valueChanges.pipe(
      distinctUntilChanged(),
      combineLatestWith(this.reload$),
      map(([sprintId, _]) => sprintId),
      filter(sprintId => sprintId != null),
      switchMap(sprintId => {
        this.loading = true;
        return this.reportService.getBurndownReport(sprintId).pipe(
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
    this.form.get('sprint')!.valueChanges.pipe(
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
    this.totalScopeDev = 0;
    this.remainingDev = 0;
    this.completedPercentDev = 0;
    this.totalScopeQA = 0;
    this.remainingQA = 0;
    this.completedPercentQA = 0;
  }

  private buildCharts(report: BurndownReport): void {
    const theme = this.themeService.theme$.value;
    const allDates = this.generateDateRange(report.start_date, report.end_date);

    // Build Dev chart
    this.buildDevChart(report, allDates, theme);

    // Build QA chart
    this.buildQAChart(report, allDates, theme);
  }

  private buildDevChart(report: BurndownReport, allDates: string[], theme: string): void {
    const dates = report.data_points.map(dp => dp.date);
    const remaining = report.data_points.map(dp => dp.remaining_dev);
    const totalScopes = report.data_points.map(dp => dp.total_scope);

    // Calculate stats
    const initialScope = totalScopes.length > 0 ? Math.max(...totalScopes) : 0;
    const lastRemaining = remaining.length > 0 ? remaining[remaining.length - 1] : 0;

    this.totalScopeDev = initialScope;
    this.remainingDev = lastRemaining;
    this.completedPercentDev = initialScope > 0
      ? Math.round(((initialScope - lastRemaining) / initialScope) * 100)
      : 0;

    // Build ideal burndown line
    const idealStep = allDates.length > 1 ? initialScope / (allDates.length - 1) : 0;
    const idealValues = allDates.map((_, i) => Math.max(0, Math.round((initialScope - idealStep * i) * 10) / 10));

    this.devGraph.data = [
      {
        x: allDates,
        y: idealValues,
        type: 'scatter',
        mode: 'lines',
        name: 'Идеальная',
        line: {
          color: '#28a745',
          width: 2,
          dash: 'dash'
        }
      },
      {
        x: dates,
        y: remaining,
        type: 'scatter',
        mode: 'lines+markers',
        name: 'Фактическая',
        line: {
          color: '#1151F3',
          width: 2
        },
        marker: {
          size: 6
        }
      }
    ];

    this.devGraph.layout = this.buildLayout(theme, 'Разработка');
  }

  private buildQAChart(report: BurndownReport, allDates: string[], theme: string): void {
    const dates = report.data_points.map(dp => dp.date);
    const remainingQA = report.data_points.map(dp => dp.remaining_qa);
    const totalScopesQA = report.data_points.map(dp => dp.total_scope_qa);

    // Calculate stats
    const initialScopeQA = totalScopesQA.length > 0 ? Math.max(...totalScopesQA) : 0;
    const lastRemainingQA = remainingQA.length > 0 ? remainingQA[remainingQA.length - 1] : 0;

    this.totalScopeQA = initialScopeQA;
    this.remainingQA = lastRemainingQA;
    this.completedPercentQA = initialScopeQA > 0
      ? Math.round(((initialScopeQA - lastRemainingQA) / initialScopeQA) * 100)
      : 0;

    // Build ideal burndown line
    const idealStep = allDates.length > 1 ? initialScopeQA / (allDates.length - 1) : 0;
    const idealValues = allDates.map((_, i) => Math.max(0, Math.round((initialScopeQA - idealStep * i) * 10) / 10));

    this.qaGraph.data = [
      {
        x: allDates,
        y: idealValues,
        type: 'scatter',
        mode: 'lines',
        name: 'Идеальная',
        line: {
          color: '#28a745',
          width: 2,
          dash: 'dash'
        }
      },
      {
        x: dates,
        y: remainingQA,
        type: 'scatter',
        mode: 'lines+markers',
        name: 'Фактическая',
        line: {
          color: '#F31155',
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

    const titleText = type === 'Разработка' ? 'Dev Burndown' : 'QA Burndown';
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

  private generateDateRange(startDate: string, endDate: string): string[] {
    const dates: string[] = [];
    const current = new Date(startDate);
    const end = new Date(endDate);

    while (current <= end) {
      dates.push(current.toLocaleDateString());
      current.setDate(current.getDate() + 1);
    }

    return dates;
  }

  reload() {
    this.reload$.next(null)
  }
}
