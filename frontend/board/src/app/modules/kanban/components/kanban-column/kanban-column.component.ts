import {Component, DestroyRef, EventEmitter, inject, Input, OnDestroy, Output} from '@angular/core';
import {KanbanColumn} from "../../models/kanban-column";
import {KanbanUser} from "../../models/kanban-user";
import { faEllipsis } from '@fortawesome/free-solid-svg-icons'
import {KanbanColumnModalService} from "../../services/kanban-column-modal.service";
import {KanbanColumnService} from "../../services/kanban-column.service";
import {takeUntil} from "rxjs/operators";
import {Subject} from "rxjs";
import {AccountService} from "../../../core/services/account.service";
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Component({
    selector: 'app-kanban-column',
    templateUrl: './kanban-column.component.html',
    styleUrls: ['./kanban-column.component.scss'],
    standalone: false
})
export class KanbanColumnComponent{
  private columnModal = inject(KanbanColumnModalService)
  private columnService = inject(KanbanColumnService)
  public accountService = inject(AccountService)
  private destroyRef = inject(DestroyRef)

  @Input() column!: KanbanColumn
  @Input() user!: KanbanUser
  @Input() search: string | null = null
  @Output() onDelete = new EventEmitter<KanbanColumn>()

  faEllipsis = faEllipsis


  openModal(e: MouseEvent) {
    this.columnModal.show(this.column).then(value => this.column = value as KanbanColumn)
    e.preventDefault();
    return false
  }

  delete(e: MouseEvent) {
    if (confirm('Удалить эту колонку?')) {
      this.columnService.delete(this.column).pipe(
        takeUntilDestroyed(this.destroyRef)
      ).subscribe(_ => {
        this.onDelete.emit(this.column)
      })
    }
    e.preventDefault()
    return false
  }
}
