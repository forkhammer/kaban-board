import {AfterViewInit, Component, DestroyRef, inject, OnInit, ViewChild} from '@angular/core';
import {KanbanUserService} from "../../services/kanban-user.service";
import {
  BehaviorSubject,
  combineLatestWith,
  distinctUntilChanged,
  filter,
  finalize,
  Observable,
  switchMap, timer
} from "rxjs";
import {KanbanUser} from "../../models/kanban-user";
import {map} from "rxjs/operators";
import {ActivatedRoute, Router} from "@angular/router";
import {FormBuilder, FormGroup} from "@angular/forms";
import { faXmark, faArrowLeft, faArrowRight, faTableList, faTableColumns, faPlus, faArrowTrendDown, faArrowTrendUp } from '@fortawesome/free-solid-svg-icons';
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
import { Sprint } from '../../models/sprint';
import { SelectModelComponent } from 'src/app/modules/ui/components/select-model/select-model.component';
import { AccountService } from 'src/app/modules/core/services/account.service';
import { SprintModalServiceService } from '../../services/sprint-modal.service';
import { UrlService } from '../../services/urls.service';

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
export class KanbanBoardComponent implements OnInit, AfterViewInit {
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
  accountService = inject(AccountService)
  sprintModal = inject(SprintModalServiceService)
  urls = inject(UrlService)

  faXmark = faXmark
  faArrowLeft = faArrowLeft
  faArrowRight = faArrowRight
  faTableList = faTableList
  faTableColumns = faTableColumns
  faArrowTrendDown = faArrowTrendDown
  faArrowTrendUp = faArrowTrendUp
  faPlus = faPlus
  COLUMN_WIDTH = 340
  KanbanView = KanbanView

  users: KanbanUser[] = []
  teams$ = new BehaviorSubject<Team[]>([])
  sprintFilter: Record<string, any>  = {}

  public isLoading = false
  public selectedUser: KanbanUser | undefined = undefined
  public searchForm: FormGroup
  public filterForm: FormGroup
  public teamId$: Observable<number | null>
  public sprintId$: Observable<number | null>
  public userId$: Observable<number | null>
  public selectedTeam: Team | null = null
  public view$: Observable<KanbanView | null>
  public selectedSprint: Sprint | null = null
  @ViewChild('sprintSelect') sprintSelect!: SelectModelComponent

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
      user: [null],
      view: [KanbanView.BOARD],
    })
    this.teamId$ = this.route.queryParams.pipe(
      map(params => params['team'] ? Number(params['team']) : null)
    );
    this.sprintId$ = this.route.queryParams.pipe(
      map(params => params['sprint'] ? Number(params['sprint']) : null)
    );
    this.userId$ = this.route.queryParams.pipe(
      map(params => params['user'] ? Number(params['user']) : null)
    );
    this.view$ = this.route.queryParams.pipe(
      map(params => params['view'] ? params['view'] : null)
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

      if (value) {
        this.sprintFilter = {team: value}
      } else {
        this.sprintFilter = {}
      }
    })

    this.teamId$.pipe(
      combineLatestWith(this.teams$.pipe(filter(value => value.length > 0))),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(([teamId, teams]) => {
      this.selectedTeam = teams.find(team => team.id === teamId) ?? teams[0]

      if (this.selectedTeam.id !== teamId) {
        this.router.navigate([], {
          queryParams: {team: this.selectedTeam.id},
          queryParamsHandling: 'merge'
        })
      }
    })

    this.sprintId$.pipe(
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(value => {
      this.filterForm.patchValue({sprint: value})
    })

    this.userId$.pipe(
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(value => {
      this.filterForm.patchValue({user: value})
    })

    this.view$.pipe(
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(value => {
      this.filterForm.patchValue({view: value})
    })

  }

  ngAfterViewInit(): void {
    this.sprintSelect.valueModel$.pipe(
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(value => {
      this.selectedSprint = value as Sprint
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
    this.filterForm.patchValue({team: team.id, sprint: null, user: null})
  }

  goToTeamBoard() {
    this.router.navigate(['/'], {queryParams: {user: null}, queryParamsHandling: 'merge'})
  }

  addSprint() {
    this.sprintModal.show()
      .then(sprint => {
        if (sprint) {
          this.filterForm.patchValue({sprint: sprint.id})
          this.sprintSelect.reload()
        }
      })
      .catch(() => {})
  }
}
