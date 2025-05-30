import {Component, Input} from '@angular/core';
import {KanbanIssue} from "../../models/kanban-issue";
import {faSquare} from "@fortawesome/free-regular-svg-icons";
import {faHammer, faBug} from "@fortawesome/free-solid-svg-icons";

@Component({
    selector: 'app-kanban-issue',
    templateUrl: './kanban-issue.component.html',
    styleUrls: ['./kanban-issue.component.scss'],
    standalone: false
})
export class KanbanIssueComponent {
  @Input() issue: KanbanIssue | null = null

  faSquare = faSquare
  faHammer = faHammer
  faBug = faBug

  get assignee() {
    return this.issue?.assignees.length ? this.issue?.assignees[0] : null
  }
}
