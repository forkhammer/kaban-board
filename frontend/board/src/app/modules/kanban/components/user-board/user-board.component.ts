import {Component, DestroyRef, ElementRef, inject, Input, ViewChild} from '@angular/core';
import {KanbanUser} from "../../models/kanban-user";
import {KanbanColumn} from "../../models/kanban-column";
import {KanbanColumnModalService} from "../../services/kanban-column-modal.service";
import {AccountService} from "../../../core/services/account.service";
import { CdkDragDrop, CdkDragRelease, CdkDragStart, moveItemInArray } from '@angular/cdk/drag-drop';
import { KanbanColumnService } from '../../services/kanban-column.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { faArrowLeft, faArrowRight, faPlus } from '@fortawesome/free-solid-svg-icons';
import { BehaviorSubject, combineLatestWith, debounceTime, distinctUntilChanged, filter, switchMap, timer } from 'rxjs';
import { Team } from '../../models/team';
import { KanbanIssue } from '../../models/kanban-issue';
import { IssueService } from '../../services/issue.service';
import { ToastService } from 'src/app/modules/core/services/toast.service';
import { environment } from 'src/environments/environment';
import { catchErrorMessages } from 'src/app/modules/core/tools/catch-error';
import { Sprint } from '../../models/sprint';
import { isEqual } from 'lodash';
import { Pagination } from 'src/app/modules/core/models/base';
import { IssueBindingService } from '../../services/issue-binding.service';

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
  private issueService = inject(IssueService)
  private issueBindingService = inject(IssueBindingService)
  private toast = inject(ToastService)

  faPlus = faPlus
  faArrowLeft = faArrowLeft
  faArrowRight = faArrowRight
  COLUMN_WIDTH = 310

  @ViewChild('UserBoardInner') userBoardInner: ElementRef | null = null
  public slidePosition = 0
  private isDrag$ = new BehaviorSubject<boolean>(false)
  public columns: KanbanColumn[] = []
  private updateColumnSignal$ = new BehaviorSubject(null)
  public user$ = new BehaviorSubject<KanbanUser | null | undefined>(null)
  public team$ = new BehaviorSubject<Team | null | undefined>(null)
  public sprint$ = new BehaviorSubject<Sprint | null | undefined>(null)
  public issues: KanbanIssue[] = []
  private timer$ = timer(0, environment.autoUpdateIssuesMin * 60 * 1000)

  @Input() set user(value : KanbanUser | undefined | null) {
    this.user$.next(value)
  }

  @Input() set team(value: Team | undefined | null) {
    this.team$.next(value)
  }

  @Input() set sprint(value: Sprint | undefined | null) {
    this.sprint$.next(value)
  }

  constructor() {
    this.updateColumnSignal$.pipe(
      switchMap(_ => this.kanbanColumnsService.list()),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(data => {
      this.columns = data as KanbanColumn[]
    })

    this.timer$.pipe(
      combineLatestWith(this.user$, this.team$, this.sprint$),
      filter(([_, user, team, sprint]) => {
        return !!team
      }),
      distinctUntilChanged(isEqual),
      debounceTime(1),
      switchMap(([_, user, team, sprint]) => {
        const query: Record<string, any> = {
          'team': team!.id,
          'limit': 1000,
        }
        if (user) {
          query['assignee'] = user.id
        }
        if (sprint) {
          query['sprint'] = sprint.id
        }

        if (sprint) {
          return this.issueBindingService.list(query).pipe(
            catchErrorMessages(this.toast)
          )
        } else {
          return this.issueService.list(query).pipe(
            catchErrorMessages(this.toast)
          )
        }
      }),
      takeUntilDestroyed()
    ).subscribe(data => {
      this.issues = (data as Pagination<KanbanIssue>).results
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
