import {Component, Input} from '@angular/core';
import {KanbanLabel} from "../../models/kanban-label";

@Component({
  selector: 'app-kanban-label',
  templateUrl: './kanban-label.component.html',
  styleUrls: ['./kanban-label.component.scss']
})
export class KanbanLabelComponent {
  @Input() label: KanbanLabel | null = null
}
