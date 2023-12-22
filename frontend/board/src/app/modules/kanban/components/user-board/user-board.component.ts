import {Component, EventEmitter, Input, OnDestroy, Output} from '@angular/core';
import {KanbanUser} from "../../models/kanban-user";
import {KanbanColumn} from "../../models/kanban-column";
import {KanbanColumnModalService} from "../../services/kanban-column-modal.service";
import {AccountService} from "../../../core/services/account.service";
import { faPlus } from '@fortawesome/free-solid-svg-icons'
import { CdkDragDrop, CdkDragEnd, CdkDragEnter, CdkDragExit, CdkDragRelease, CdkDragStart, moveItemInArray } from '@angular/cdk/drag-drop';
import { KanbanColumnService } from '../../services/kanban-column.service';
import { Subject, takeUntil } from 'rxjs';

@Component({
  selector: 'app-user-board',
  templateUrl: './user-board.component.html',
  styleUrls: ['./user-board.component.scss']
})
export class UserBoardComponent implements OnDestroy {
  faPlus = faPlus

  @Input() user!: KanbanUser
  @Input() columns: KanbanColumn[] = []
  @Input() columnWidth = 350
  @Input() search: string | null = null
  @Output() onAddColumn: EventEmitter<KanbanColumn> = new EventEmitter<KanbanColumn>()
  @Output() onDeleteColumn: EventEmitter<KanbanColumn> = new EventEmitter<KanbanColumn>()
  @Output() onUpdateColumns: EventEmitter<KanbanColumn[]> = new EventEmitter<KanbanColumn[]>()
  @Output() onDrag: EventEmitter<boolean> = new EventEmitter<boolean>()
  private destroy$ = new Subject()

  constructor(
    private columnModal: KanbanColumnModalService,
    public accountService: AccountService,
    private columnService: KanbanColumnService,
  ) {
  }

  ngOnDestroy(): void {
      this.destroy$.next(null)
      this.destroy$.complete()
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

  dropColumn(e: CdkDragDrop<KanbanColumn[]>) {
    moveItemInArray(this.columns, e.previousIndex, e.currentIndex)
    this.columns.map((column, index) => {
      column.order = index
    })
    this.columnService.saveOrdering(this.columns).pipe(
      takeUntil(this.destroy$)
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
