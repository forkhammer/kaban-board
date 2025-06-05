import {Component, DestroyRef, ElementRef, inject, OnInit, ViewChild} from '@angular/core';
import {KanbanUserService} from "../../services/kanban-user.service";
import {
  BehaviorSubject,
  combineLatestWith,
  distinctUntilChanged,
  filter,
  finalize,
  Observable,
  of,
  switchMap, timer
} from "rxjs";
import {KanbanUser} from "../../models/kanban-user";
import {map} from "rxjs/operators";
import {KanbanColumn} from "../../models/kanban-column";
import {ActivatedRoute, Router} from "@angular/router";
import {FormBuilder, FormGroup} from "@angular/forms";
import { faXmark, faArrowLeft, faArrowRight, faTableList, faTableColumns, faPlus } from '@fortawesome/free-solid-svg-icons';
import {TitleService} from "../../../core/services/title.service";
import {TeamService} from "../../services/team.service";
import {GitlabSyncService} from "../../services/gitlab-sync.service";
import {environment} from "../../../../../environments/environment";
import { Team } from '../../models/team';
import { Group } from '../../models/group';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { catchErrorMessages } from 'src/app/modules/core/tools/catch-error';
import { ToastService } from 'src/app/modules/core/services/toast.service';
import { SprintService } from '../../services/sprint.service';
import { isEqual } from 'lodash';

enum KanbanView {
  LIST = 'list',
  BOARD = 'board'
}

@Component({
    selector: 'app-kanban-board',
    templateUrl: './kanban-board.component.html',
    styleUrls: ['./kanban-board.component.scss'],
    standalone: false
})
export class KanbanBoardComponent implements OnInit {
  private kanbanUserService = inject(KanbanUserService)

  private route = inject(ActivatedRoute)
  private router = inject(Router)
  private builder = inject(FormBuilder)
  private title = inject(TitleService)
  public teamService = inject(TeamService)
  private syncService = inject(GitlabSyncService)
  private destroyRef = inject(DestroyRef)
  private toast = inject(ToastService)
  sprintService = inject(SprintService)

  faXmark = faXmark
  faArrowLeft = faArrowLeft
  faArrowRight = faArrowRight
  faTableList = faTableList
  faTableColumns = faTableColumns
  faPlus = faPlus
  COLUMN_WIDTH = 340
  KanbanView = KanbanView

  users: KanbanUser[] = []
  teams$ = new BehaviorSubject<Team[]>([])

  public isLoading = false
  public selectedUser: KanbanUser | undefined = undefined
  public searchForm: FormGroup
  public filterForm: FormGroup
  public teamId$: Observable<number | null>
  public selectedTeam: Team | null = null
  public search$: Observable<string | null>

  public otherGroup: Group = {
    id: 0,
    title: 'Остальные',
  }

  get view(): KanbanView {
    return this.filterForm.get('view')?.value ?this.filterForm.get('view')?.value as KanbanView : KanbanView.BOARD
  }

  constructor() {
    this.searchForm = this.builder.group({
      search: [''],
    })
    this.filterForm = this.builder.group({
      team: [null],
      sprint: [null],
      view: [KanbanView.BOARD],
    })
    this.teamId$ = this.route.queryParams.pipe(
      map(params => params['team'] ? Number(params['team']) : null)
    );
    this.search$ = this.route.queryParams.pipe(
      map(params => params['search'] ? params['search'] : null)
    );

    this.teamService.list().pipe(
      catchErrorMessages(this.toast),
      takeUntilDestroyed(),
    ).subscribe(data => {
      this.teams$.next(data as Team[])
    })
  }

  ngOnInit() {
    this.isLoading = true
    this.title.setTitle('General board')

    const userId$ = this.route.queryParams.pipe(
      map(params => Number(params['user']))
    );

    const users$ = timer(0, environment.autoUpdateIssuesMin * 60 * 1000).pipe(
      switchMap(_ => this.kanbanUserService.listUsers().pipe(
        catchErrorMessages(this.toast),
        finalize(() => this.isLoading = false),
      )),
      takeUntilDestroyed(this.destroyRef)
    )

    userId$.pipe(
      combineLatestWith(users$),
      takeUntilDestroyed(this.destroyRef),
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

    this.filterForm?.valueChanges.pipe(
      distinctUntilChanged(isEqual),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(value => {
      this.router.navigate(['/'], {queryParams:value, queryParamsHandling: 'merge'})
    })

    this.teamId$.pipe(
      filter(value => value != undefined),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(value => {
      this.filterForm.patchValue({team: value})
    })

    this.teamId$.pipe(
      combineLatestWith(this.teams$.pipe(filter(value => value.length > 0))),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(([teamId, teams]) => {
      this.selectedTeam = teams.find(team => team.id === teamId) ?? teams[0]
    })
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

  viewBoard() {
    this.filterForm.patchValue({view: KanbanView.BOARD})
  }

  viewList() {
    this.filterForm.patchValue({view: KanbanView.LIST})
  }

  getTeamTitle(team: Team) {
    return team.title.slice(0, 2)
  }

  selectTeam(team: Team) {
    this.filterForm.patchValue({team: team.id})
  }

  goToTeamBoard() {
    this.router.navigate(['/'], {queryParams: {user: null}, queryParamsHandling: 'merge'})
  }
}
