import { Component, EventEmitter, Input, Output } from '@angular/core';
import { Sprint } from '../../models/sprint';

@Component({
  selector: 'app-sprint-list',
  standalone: false,
  templateUrl: './sprint-list.component.html',
  styleUrl: './sprint-list.component.scss'
})
export class SprintListComponent {
  @Input() sprints: Sprint[] = []
  @Output() onDelete = new EventEmitter<Sprint>()
}
