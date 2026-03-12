import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import * as PlotlyJS from 'plotly.js-dist-min';
import { PlotlyModule } from 'angular-plotly.js';

PlotlyModule.plotlyjs = PlotlyJS;

@Component({
  selector: 'app-burndown-report-page',
  standalone: true,
  imports: [CommonModule, PlotlyModule],
  templateUrl: './burndown-report-page.component.html',
  styleUrl: './burndown-report-page.component.scss'
})
export class BurndownReportPageComponent {
  // Данные для burndown диаграммы
  public graph = {
    data: [
      // Идеальная линия сгорания
      {
        x: ['Day 1', 'Day 2', 'Day 3', 'Day 4', 'Day 5', 'Day 6', 'Day 7', 'Day 8', 'Day 9', 'Day 10', 'Day 11', 'Day 12', 'Day 13', 'Day 14'],
        y: [100, 92.8, 85.7, 78.5, 71.4, 64.2, 57.1, 50, 42.8, 35.7, 28.5, 21.4, 14.2, 7.1, 0],
        type: 'scatter',
        mode: 'lines+markers',
        name: 'Ideal',
        line: {
          color: '#28a745',
          width: 2,
          dash: 'dash'
        },
        marker: {
          size: 6
        }
      },
      // Фактическая линия
      {
        x: ['Day 1', 'Day 2', 'Day 3', 'Day 4', 'Day 5', 'Day 6', 'Day 7', 'Day 8', 'Day 9', 'Day 10', 'Day 11', 'Day 12', 'Day 13', 'Day 14'],
        y: [100, 95, 88, 82, 78, 72, 68, 65, 58, 52, 48, 42, 35, 28, 20],
        type: 'scatter',
        mode: 'lines+markers',
        name: 'Actual',
        line: {
          color: '#dc3545',
          width: 2
        },
        marker: {
          size: 6
        }
      }
    ],
    layout: {
      title: {
        text: 'Sprint Burndown Chart',
        font: {
          size: 20
        }
      },
      xaxis: {
        title: 'Sprint Days',
        showgrid: true,
        gridcolor: '#e0e0e0'
      },
      yaxis: {
        title: 'Remaining Story Points',
        showgrid: true,
        gridcolor: '#e0e0e0',
        range: [0, 110]
      },
      legend: {
        x: 0.02,
        y: 0.98,
        bgcolor: 'rgba(255, 255, 255, 0.8)',
        bordercolor: '#e0e0e0',
        borderwidth: 1
      },
      plot_bgcolor: '#fafafa',
      paper_bgcolor: '#ffffff',
      margin: {
        l: 60,
        r: 40,
        t: 60,
        b: 60
      }
    },
    config: {
      responsive: true,
      displayModeBar: true,
      displaylogo: false
    }
  };
}
