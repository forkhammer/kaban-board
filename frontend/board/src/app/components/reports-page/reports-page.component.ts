import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { IconDefinition, FaIconComponent } from '@fortawesome/angular-fontawesome';
import {faArrowTrendDown, faArrowTrendUp, faChartLine} from '@fortawesome/free-solid-svg-icons';

interface ReportCard {
  id: string;
  title: string;
  description: string;
  route: string;
  icon: IconDefinition
}

@Component({
  selector: 'app-reports-page',
  standalone: true,
  imports: [CommonModule, RouterLink, FaIconComponent],
  templateUrl: './reports-page.component.html',
  styleUrl: './reports-page.component.scss'
})
export class ReportsPageComponent {
  reports: ReportCard[] = [
    {
      id: 'burndown',
      title: 'Burndown Chart',
      description: 'График сгорания работ - показывает оставшийся объем работ по времени, помогает отслеживать прогресс спринта',
      route: '/reports/burndown',
      icon: faArrowTrendDown,
    },
    {
      id: 'burnup',
      title: 'Burnup Chart',
      description: 'График нарастания работ - отображает выполненную работу и общий объем по времени, показывает изменения в объеме',
      route: '/reports/burnup',
      icon: faArrowTrendUp
    },
    {
      id: 'wip',
      title: 'WIP Chart',
      description: 'Work In Progress - показывает количество задач одновременно находящихся в работе по времени',
      route: '/reports/wip',
      icon: faChartLine
    }
  ];
}
