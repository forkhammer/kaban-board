import {Component, ElementRef, OnDestroy, OnInit, ViewChild} from '@angular/core';
import {KanbanUserService} from "../../services/kanban-user.service";
import {
  BehaviorSubject,
  combineLatestWith,
  distinctUntilChanged,
  filter,
  finalize,
  Observable,
  of,
  Subject,
  switchMap, timer
} from "rxjs";
import {KanbanUser} from "../../models/kanban-user";
import {map, takeUntil} from "rxjs/operators";
import {KanbanColumn} from "../../models/kanban-column";
import {KanbanColumnService} from "../../services/kanban-column.service";
import {ActivatedRoute, Router} from "@angular/router";
import {FormBuilder, FormGroup} from "@angular/forms";
import { faXmark, faArrowLeft, faArrowRight } from '@fortawesome/free-solid-svg-icons';
import {TitleService} from "../../../core/services/title.service";
import {TeamService} from "../../services/team.service";
import {GitlabSyncService} from "../../services/gitlab-sync.service";
import {environment} from "../../../../../environments/environment";
import { Team } from '../../models/team';

@Component({
  selector: 'app-kanban-board',
  templateUrl: './kanban-board.component.html',
  styleUrls: ['./kanban-board.component.scss']
})
export class KanbanBoardComponent implements OnInit, OnDestroy {
  faXmark = faXmark
  faArrowLeft = faArrowLeft
  faArrowRight = faArrowRight
  COLUMN_WIDTH = 340

  private destroy$ = new Subject()
  public users: KanbanUser[] = []
  public columns: KanbanColumn[] = []
  public isLoading = false
  public selectedUser: KanbanUser | undefined = undefined
  public searchForm: FormGroup
  public filterForm: FormGroup
  public teamId$: Observable<number | null>
  public team: Team | null = null
  public slidePosition = 0
  @ViewChild('UserBoardInner') userBoardInner: ElementRef | null = null
  public search$: Observable<string | null>
  private updateColumnSignal$ = new BehaviorSubject(null)
  private isDrag$ = new BehaviorSubject<boolean>(false)

  constructor(
    private kanbanUserService: KanbanUserService,
    private kanbanColumnsService: KanbanColumnService,
    private route: ActivatedRoute,
    private router: Router,
    private builder: FormBuilder,
    private title: TitleService,
    public teamService: TeamService,
    private syncService: GitlabSyncService,
  ) {
    this.searchForm = this.builder.group({
      search: [''],
    })
    this.filterForm = this.builder.group({
      team: [null]
    })
    this.teamId$ = this.route.queryParams.pipe(
      map(params => params['team'] ? Number(params['team']) : null)
    );
    this.search$ = this.route.queryParams.pipe(
      map(params => params['search'] ? params['search'] : null)
    );
  }

  ngOnInit() {
    this.isLoading = true
    this.title.setTitle('General board')

    const userId$ = this.route.queryParams.pipe(
      map(params => Number(params['user']))
    );

    const users$ = timer(0, environment.autoUpdateIssuesMin * 60 * 1000).pipe(
      switchMap(_ => this.kanbanUserService.listUsers().pipe(
        finalize(() => this.isLoading = false),
      )),
      takeUntil(this.destroy$)
    )

    userId$.pipe(
      combineLatestWith(users$),
      takeUntil(this.destroy$)
    ).subscribe(([userId, resp]) => {
      this.users = resp.users
      this.syncService.updateTime$.next(resp.updateTime ? new Date(resp.updateTime) : null)

      if (userId) {
        const user = this.getUserById(userId)
        this.selectUser(user)
        this.title.setTitle(`${user?.name} board`)
      } else {
        this.selectUser(undefined)
        this.title.setTitle('General board')
      }
    })

    this.updateColumnSignal$.pipe(
      switchMap(_ => this.kanbanColumnsService.list()),
      takeUntil(this.destroy$)
    ).subscribe(data => {
      this.columns = data as KanbanColumn[]
    })

    this.filterForm.get('team')?.valueChanges.pipe(
      distinctUntilChanged(),
      takeUntil(this.destroy$)
    ).subscribe(value => {
      this.router.navigate(['/'], {queryParams:{team: value ? value : ''}, queryParamsHandling: 'merge'})
    })

    this.teamId$.pipe(
      filter(value => value != undefined),
      takeUntil(this.destroy$)
    ).subscribe(value => {
      this.filterForm.patchValue({team: value})
    })

    this.teamId$.pipe(
      switchMap(value => value ? this.teamService.get(value) : of(null)),
      takeUntil(this.destroy$)
    ).subscribe(data => {
      this.team = data
    })
  }

  ngOnDestroy() {
    this.destroy$.next(null)
    this.destroy$.complete()
  }

  goToUserBoard(user: KanbanUser) {
    this.router.navigate(['/'], {queryParams: {user: user.id}, queryParamsHandling: 'merge'})
  }

  selectUser(user: KanbanUser | undefined) {
    this.selectedUser = user
  }

  trackByUser(index: number, user: KanbanUser) {
    return user.id
  }

  getUserById(id: number): KanbanUser | undefined {
    return this.users.find(user => user.id === id)
  }

  clearSearch(e: MouseEvent) {
    this.searchForm.patchValue({search:''})
    e.preventDefault()
    return false
  }

  catchAddColumn(column: KanbanColumn) {
    this.columns.push(column)
  }

  catchDeleteColumn(column: KanbanColumn) {
    this.columns.splice(this.columns.findIndex(c => c.id == column.id), 1)
  }

  catchUpdateColumns(columns: KanbanColumn[]) {
    this.updateColumnSignal$.next(null)
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

  getScreenColumnsCount() {
    return Math.floor(this.userBoardInner?.nativeElement?.offsetWidth / this.COLUMN_WIDTH)
  }

  getSlideStep() {
    return this.getScreenColumnsCount()
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

  catchDrag(e: boolean) {
    this.isDrag$.next(e)
  }
}
