import {Component, DestroyRef, ElementRef, EventEmitter, inject, Input, Output, ViewChild} from '@angular/core';
import {KanbanUser} from "../../models/kanban-user";
import {KanbanColumn} from "../../models/kanban-column";
import {KanbanColumnModalService} from "../../services/kanban-column-modal.service";
import {AccountService} from "../../../core/services/account.service";
import { CdkDragDrop, CdkDragRelease, CdkDragStart, moveItemInArray } from '@angular/cdk/drag-drop';
import { KanbanColumnService } from '../../services/kanban-column.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { faArrowLeft, faArrowRight, faPlus } from '@fortawesome/free-solid-svg-icons';
import { BehaviorSubject, switchMap } from 'rxjs';
import { Team } from '../../models/team';

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
  private kanbanColumnsService = inject(KanbanColumnService)

  faPlus = faPlus
  faArrowLeft = faArrowLeft
  faArrowRight = faArrowRight
  COLUMN_WIDTH = 340

  @Input() search: string | null = null
  @ViewChild('UserBoardInner') userBoardInner: ElementRef | null = null
  public slidePosition = 0
  private isDrag$ = new BehaviorSubject<boolean>(false)
  public columns: KanbanColumn[] = []
  private updateColumnSignal$ = new BehaviorSubject(null)
  public user$ = new BehaviorSubject<KanbanUser | null>(null)
  public team$ = new BehaviorSubject<Team | null>(null)

  @Input() set user(value : KanbanUser | null) {
    this.user$.next(value)
  }

  @Input() set team(value: Team | null) {
    this.team$.next(value)
  }

  constructor() {
    this.updateColumnSignal$.pipe(
      switchMap(_ => this.kanbanColumnsService.list()),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(data => {
      this.columns = data as KanbanColumn[]
    })
  }

  trackByColumn(index: number, column: KanbanColumn) {
    return column.id
  }

  addColumn(e: MouseEvent) {
    this.columnModal.show(null).then(value => {
      if (value) {
        this.columns.push(value as KanbanColumn)
      }
    })
    e.preventDefault()
    return false
  }

  catchDeleteColumn(column: KanbanColumn) {
    this.columns.splice(this.columns.findIndex(c => c.id == column.id), 1)
  }

  dropColumn(e: CdkDragDrop<KanbanColumn[]>) {
    moveItemInArray(this.columns, e.previousIndex, e.currentIndex)
    this.columns.map((column, index) => {
      column.order = index
    })
    this.columnService.saveOrdering(this.columns).pipe(
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(items => {
      this.updateColumnSignal$.next(null)
    })
  }

  dragStart(e: CdkDragStart<any>) {
    this.isDrag$.next(true)
  }

  dragEnd(e: CdkDragRelease<any>) {
    this.isDrag$.next(false)
  }

  swipeLeft(e: Event) {
    if (!this.isDrag$.value) {
      this.slideRight(null)
    }
  }

  swipeRight(e: Event) {
    if (!this.isDrag$.value) {
      this.slideLeft(null)
    }
  }

  slideLeft(e: MouseEvent | null) {
    if (this.slidePosition + this.getSlideStep() > 0) {
      this.slidePosition = 0
    } else {
      this.slidePosition += this.getSlideStep()
    }
    return false
  }

  slideRight(e: MouseEvent | null) {
    this.slidePosition -= this.getSlideStep()
    return false
  }

  getSlideStep() {
    return this.getScreenColumnsCount()
  }

  getScreenColumnsCount() {
    return Math.floor(this.userBoardInner?.nativeElement?.offsetWidth / this.COLUMN_WIDTH)
  }

  getUserBoardStyles(): {[p:string]: any} {
    return {
      'transform': `translateX(${this.slidePosition * this.COLUMN_WIDTH}px)`,
    }
  }

  getScreenStartColumn() {
    return -this.slidePosition
  }

  getScreenEndColumn() {
    return this.getScreenColumnsCount() + this.getScreenStartColumn()
  }

  getActiveColumns(teamId: number | null): KanbanColumn[] {
    let columns = this.filterColumnByTeam(teamId)
    if (columns.length === 0) {
      columns = this.filterColumnByTeam(null)
    }
    return columns
  }

  filterColumnByTeam(teamId: number | null): KanbanColumn[] {
    return this.columns.filter(column => {
      if (teamId) {
        return column.team_id === teamId
      } else {
        return column.team_id === null
      }
    })
  }

  catchDrag(e: boolean) {
    this.isDrag$.next(e)
  }
}
