import {Component, DestroyRef, EventEmitter, inject, Input, Output} from '@angular/core';
import {KanbanUser} from "../../models/kanban-user";
import {KanbanColumn} from "../../models/kanban-column";
import {KanbanColumnModalService} from "../../services/kanban-column-modal.service";
import {AccountService} from "../../../core/services/account.service";
import { faPlus } from '@fortawesome/free-solid-svg-icons'
import { CdkDragDrop, CdkDragRelease, CdkDragStart, moveItemInArray } from '@angular/cdk/drag-drop';
import { KanbanColumnService } from '../../services/kanban-column.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Component({
    selector: 'app-user-board',
    templateUrl: './user-board.component.html',
    styleUrls: ['./user-board.component.scss'],
    standalone: false
})
export class UserBoardComponent {
  private columnModal = inject(KanbanColumnModalService)
  public accountService = inject(AccountService)
  private columnService = inject(KanbanColumnService)
  private destroyRef = inject(DestroyRef)

  faPlus = faPlus

  @Input() user!: KanbanUser
  @Input() columns: KanbanColumn[] = []
  @Input() columnWidth = 350
  @Input() search: string | null = null
  @Output() onAddColumn: EventEmitter<KanbanColumn> = new EventEmitter<KanbanColumn>()
  @Output() onDeleteColumn: EventEmitter<KanbanColumn> = new EventEmitter<KanbanColumn>()
  @Output() onUpdateColumns: EventEmitter<KanbanColumn[]> = new EventEmitter<KanbanColumn[]>()
  @Output() onDrag: EventEmitter<boolean> = new EventEmitter<boolean>()

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

  dropColumn(e: CdkDragDrop<KanbanColumn[]>) {
    moveItemInArray(this.columns, e.previousIndex, e.currentIndex)
    this.columns.map((column, index) => {
      column.order = index
    })
    this.columnService.saveOrdering(this.columns).pipe(
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(items => {
      this.onUpdateColumns.emit(items)
    })
  }

  dragStart(e: CdkDragStart<any>) {
    this.onDrag.emit(true)
  }

  dragEnd(e: CdkDragRelease<any>) {
    this.onDrag.emit(false)
  }
}
