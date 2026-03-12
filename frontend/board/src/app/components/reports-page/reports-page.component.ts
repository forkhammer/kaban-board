import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';

interface ReportCard {
  id: string;
  title: string;
  description: string;
  route: string;
}

@Component({
  selector: 'app-reports-page',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './reports-page.component.html',
  styleUrl: './reports-page.component.scss'
})
export class ReportsPageComponent {
  reports: ReportCard[] = [
    {
      id: 'burndown',
      title: 'Burndown Chart',
      description: 'График сгорания работ - показывает оставшийся объем работ по времени, помогает отслеживать прогресс спринта',
      route: '/reports/burndown'
    },
    {
      id: 'burnup',
      title: 'Burnup Chart',
      description: 'График нарастания работ - отображает выполненную работу и общий объем по времени, показывает изменения в объеме',
      route: '/reports/burnup'
    }
  ];
}
