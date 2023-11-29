import {Component, EventEmitter, Input, Output} from '@angular/core';
import {KanbanUser} from "../../models/kanban-user";
import {KanbanColumn} from "../../models/kanban-column";
import {KanbanColumnModalService} from "../../services/kanban-column-modal.service";
import {AccountService} from "../../../core/services/account.service";
import { faPlus } from '@fortawesome/free-solid-svg-icons'

@Component({
  selector: 'app-user-board',
  templateUrl: './user-board.component.html',
  styleUrls: ['./user-board.component.scss']
})
export class UserBoardComponent {
  faPlus = faPlus

  @Input() user!: KanbanUser
  @Input() columns: KanbanColumn[] = []
  @Input() columnWidth = 350
  @Input() search: string | null = null
  @Output() onAddColumn: EventEmitter<KanbanColumn> = new EventEmitter<KanbanColumn>()
  @Output() onDeleteColumn: EventEmitter<KanbanColumn> = new EventEmitter<KanbanColumn>()

  constructor(private columnModal: KanbanColumnModalService, public accountService: AccountService) {
  }

  trackByColumn(index: number, column: KanbanColumn) {
    return column.id
  }

  addColumn(e: MouseEvent) {
    this.columnModal.show(null).then(value => {
      if (value) {
        this.onAddColumn.emit(value as KanbanColumn)
      }
    })
    e.preventDefault()
    return false
  }

  catchDeleteColumn(column: KanbanColumn) {
    this.onDeleteColumn.emit(column)
  }
}
