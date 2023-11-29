import {Component, EventEmitter, Input, OnDestroy, Output} from '@angular/core';
import {KanbanColumn} from "../../models/kanban-column";
import {KanbanUser} from "../../models/kanban-user";
import { faEllipsis } from '@fortawesome/free-solid-svg-icons'
import {KanbanColumnModalService} from "../../services/kanban-column-modal.service";
import {KanbanColumnService} from "../../services/kanban-column.service";
import {takeUntil} from "rxjs/operators";
import {Subject} from "rxjs";
import {AccountService} from "../../../core/services/account.service";

@Component({
  selector: 'app-kanban-column',
  templateUrl: './kanban-column.component.html',
  styleUrls: ['./kanban-column.component.scss']
})
export class KanbanColumnComponent implements OnDestroy{
  @Input() column!: KanbanColumn
  @Input() user!: KanbanUser
  @Input() search: string | null = null
  @Output() onDelete = new EventEmitter<KanbanColumn>()

  faEllipsis = faEllipsis

  private destroy$ = new Subject()

  constructor(
    private columnModal: KanbanColumnModalService,
    private columnService: KanbanColumnService,
    public accountService: AccountService
  ) {
  }

  ngOnDestroy() {
    this.destroy$.next(null)
    this.destroy$.complete()
  }

  openModal(e: MouseEvent) {
    this.columnModal.show(this.column).then(value => this.column = value as KanbanColumn)
    e.preventDefault();
    return false
  }

  delete(e: MouseEvent) {
    if (confirm('Удалить эту колонку?')) {
      this.columnService.delete(this.column).pipe(
        takeUntil(this.destroy$)
      ).subscribe(_ => {
        this.onDelete.emit(this.column)
      })
    }
    e.preventDefault()
    return false
  }
}
