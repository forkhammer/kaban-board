import { Component, DestroyRef, EventEmitter, inject, Input, NgZone, Output } from '@angular/core';
import { BehaviorSubject, catchError, combineLatestWith, debounceTime, distinctUntilChanged, EMPTY, filter, merge, of, Subject, switchMap, timer } from 'rxjs';
import { KanbanUser } from '../../models/kanban-user';
import { Team } from '../../models/team';
import { Sprint } from '../../models/sprint';
import { environment } from 'src/environments/environment';
import { isEqual } from 'lodash';
import { catchErrorMessages } from 'src/app/modules/core/tools/catch-error';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { IssueGroup, KanbanIssue } from '../../models/kanban-issue';
import { IssueService } from '../../services/issue.service';
import { ToastService } from 'src/app/modules/core/services/toast.service';
import { faPlus } from '@fortawesome/free-solid-svg-icons'
import { Pagination } from 'src/app/modules/core/models/base';
import { AccountService } from 'src/app/modules/core/services/account.service';
import { IssueBindingService } from '../../services/issue-binding.service';
import { BindIssueModalService } from '../../services/bind-issue-modal.service';
import { CdkDragDrop, moveItemInArray } from '@angular/cdk/drag-drop';
import { IssueBindingModalService } from '../../services/issue-binding-modal.service';
import { BindingDeletedEventData, BindingUpdatedEventData, SprintWebsocketService } from '../../services/sprint-websocket.service';

@Component({
  selector: 'app-issue-table',
  standalone: false,
  templateUrl: './issue-table.component.html',
  styleUrl: './issue-table.component.scss'
})
export class IssueTableComponent {
  private issueService = inject(IssueService)
  private issueBindingService = inject(IssueBindingService)
  private toast = inject(ToastService)
  public accountService = inject(AccountService)
  private bindIssueModal = inject(BindIssueModalService)
  private destroyRef = inject(DestroyRef)
  private issueBindingModal = inject(IssueBindingModalService)
  private sprintWs = inject(SprintWebsocketService)
  private zone = inject(NgZone)

  faPlus = faPlus

  @Output() sprintStateChanged = new EventEmitter<void>()

  public user$ = new BehaviorSubject<KanbanUser | null | undefined>(null)
  public team$ = new BehaviorSubject<Team | null | undefined>(null)
  public sprint$ = new BehaviorSubject<Sprint | null | undefined>(null)
  private timer$ = timer(0, environment.autoUpdateIssuesMin * 60 * 1000)
  private wsReload$ = new BehaviorSubject<number>(0)
  public issues: KanbanIssue[] = []
  public groupedIssues: IssueGroup[] = []
  public issuePage: Pagination<KanbanIssue> | null = null
  public isLoading = false
  isPageMore$ = new BehaviorSubject<boolean>(false);
  isLoadMore$ = new BehaviorSubject<boolean>(false);
  highlightedIds = new Set<number>()
  highlightedTimeout = 1000

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
    this.isPageMore$.pipe(
      filter(Boolean),
      takeUntilDestroyed()
    ).subscribe(data => {
      this.isLoadMore$.next(data)
    })

    this.sprint$.pipe(
      distinctUntilChanged((a, b) => a?.id === b?.id),
      takeUntilDestroyed()
    ).subscribe(sprint => {
      if (sprint) {
        this.sprintWs.connect(sprint.id)
      } else {
        this.sprintWs.disconnect()
      }
    })

    this.sprintWs.events$.pipe(
      takeUntilDestroyed()
    ).subscribe(event => {
      switch (event.type) {
        case 'binding_updated':  this.onBindingUpdated(event.data as BindingUpdatedEventData); break
        case 'binding_deleted':  this.onBindingDeleted(event.data as BindingDeletedEventData); break
        case 'binding_created': this.onBindingCreated(event.data as BindingUpdatedEventData); break
        case 'binding_ordering': this.reloadIssues(); break
      }
    })

    this.destroyRef.onDestroy(() => this.sprintWs.disconnect())

    merge(this.timer$, this.wsReload$).pipe(
      combineLatestWith(this.user$, this.team$, this.sprint$),
      debounceTime(100),
      filter(([_, user, team, sprint]) => {
        return !!team
      }),
      distinctUntilChanged(isEqual),
      debounceTime(1),
      combineLatestWith(this.isLoadMore$),
      switchMap(([data, isLoadingMore]) => {
        const [_, user, team, sprint] = data as [any, KanbanUser, Team, Sprint]
        const query: Record<string, any> = {
          'team': team!.id,
          'limit': sprint ? 10000 : 50,
          'page': this.isPageMore$.value ? this.issuePage!.page + 1 : 1
        }
        if (user) {
          query['assignee'] = user.id
        }
        if (sprint) {
          query['sprint'] = sprint.id
        }
        this.isLoading = true

        if (sprint) {
          return this.issueBindingService.list(query).pipe(
            combineLatestWith(of(this.isPageMore$.value)),
            catchErrorMessages(this.toast, () => this.isLoading = false)
          )
        } else {
          return this.issueService.list(query).pipe(
            combineLatestWith(of(this.isPageMore$.value)),
            catchErrorMessages(this.toast, () => this.isLoading = false)
          )
        }

      }),
      takeUntilDestroyed()
    ).subscribe(([data, isLoadingMore]) => {
      this.isLoading = false

      if (isLoadingMore && this.issuePage) {
        this.issuePage = (data as Pagination<KanbanIssue>)
        this.setIssues(this.issues.concat(this.issuePage.results))
        this.isPageMore$.next(false)
      } else {
        this.issuePage = (data as Pagination<KanbanIssue>)
        this.setIssues([...this.issuePage.results].sort((a, b) => {
          const oa = a.order ?? ''
          const ob = b.order ?? ''
          if (oa === '' && ob === '') return 0
          if (oa === '') return 1
          if (ob === '') return -1
          return oa.localeCompare(ob)
        }))
      }
    })
  }

  private onBindingUpdated(data: BindingUpdatedEventData): void {
    this.sprintStateChanged.emit()
    this.groupedIssues = this.getGroupedIssues()
    if (data.account_id === this.accountService.user$.value?.id) return

    const idx = this.issues.findIndex(i => i.bindingId === data.id)
    if (idx === -1) return

    this.issueBindingService.get(data.id).pipe(
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(updated => {
      this.issues[idx] = updated
      this.groupedIssues = this.getGroupedIssues()
      this.addHighlighted(data.id)
    })
  }

  private onBindingDeleted(data: BindingDeletedEventData): void {
    this.setIssues(this.issues.filter(i => i.bindingId !== data.id))
    this.sprintStateChanged.emit()
    this.groupedIssues = this.getGroupedIssues()
  }

  private reloadIssues(): void {
    this.wsReload$.next(this.wsReload$.value + 1)
  }

  private onBindingCreated(data: BindingUpdatedEventData): void {
    this.wsReload$.next(this.wsReload$.value + 1)
    this.sprintStateChanged.emit()
    this.groupedIssues = this.getGroupedIssues()
    if (data.account_id !== this.accountService.user$.value?.id) {
      this.addHighlighted(data.id)
    }
  }

  appendIssue() {
    this.bindIssueModal.show().then(issue => {
      if (issue && this.sprint$.value) {
        this.issueService.bindToSprint((issue as KanbanIssue).id, this.sprint$.value.id, this.user$.value?.id ?? null).pipe(
          catchErrorMessages(this.toast),
          takeUntilDestroyed(this.destroyRef),
        ).subscribe(data => {
          this.issues.push(data)
          this.groupedIssues = this.getGroupedIssues()
          this.sprintStateChanged.emit()
        })
      }
    }, () => {})
  }

  createIssue() {
    this.issueBindingModal.show({
      assigneeId: this.user$.value?.id ?? null,
      sprintId: this.sprint$.value?.id ?? null
    }).then(
      (result) => {
        if (result) {
          this.issues.push(result)
          this.groupedIssues = this.getGroupedIssues()
          this.sprintStateChanged.emit()
        }
      },
      (err) => {}
    )
  }

  unbindIssue(bindingId: number) {
    this.setIssues(this.issues.filter(issue => issue.bindingId !== bindingId))
    this.sprintStateChanged.emit()
  }

  onSaveIssue(issue: KanbanIssue) {
    this.issueBindingService.save(issue).pipe(
      catchErrorMessages(this.toast),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(data => {
      const idx = this.issues.findIndex(i => i.bindingId === issue.bindingId)
      if (idx !== -1) {
        Object.assign(this.issues[idx], data)
      }
      this.groupedIssues = this.getGroupedIssues()
    })
  }

  onDeleteIssue(issue: KanbanIssue) {
    this.issueBindingService.delete(issue).pipe(
      catchErrorMessages(this.toast),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(() => {
      this.unbindIssue(issue.bindingId!)
    })
  }

  dropIssue(event: CdkDragDrop<KanbanIssue[]>): void {
    const flat = this.groupedIssues.flatMap(g => g.issues)
    moveItemInArray(flat, event.previousIndex, event.currentIndex)
    this.recalculateOrder(flat)
    this.setIssues(flat)
    this.issueBindingService.saveOrdering(flat).pipe(
      catchErrorMessages(this.toast),
      takeUntilDestroyed(this.destroyRef),
    ).subscribe()
  }

  private recalculateOrder(issues: KanbanIssue[]): void {
    const grouped = new Map<number, KanbanIssue[]>()
    for (const issue of issues) {
      const groupId = issue.assignee?.id ?? 0
      if (!grouped.has(groupId)) {
        grouped.set(groupId, [])
      }
      grouped.get(groupId)!.push(issue)
    }
    for (const [groupId, items] of grouped) {
      items.forEach((issue, index) => {
        issue.order = `${groupId}${String(index).padStart(8, '0')}`
      })
    }
  }

  setIssues(issues: KanbanIssue[]): void {
    this.issues = issues
    this.groupedIssues = this.getGroupedIssues()
  }

  getGroupedIssues(): IssueGroup[] {
    const map = new Map<number, IssueGroup>()
    const NO_GROUP_ID = -1

    for (const issue of this.issues) {
      const group = issue.assignee?.groups?.[0]
      const key = group?.id ?? NO_GROUP_ID
      if (!map.has(key)) {
        map.set(key, { groupId: key, groupTitle: group?.title ?? 'Без группы', issues: [] })
      }
      map.get(key)!.issues.push(issue)
    }

    return Array.from(map.values()).sort((a, b) => a.groupTitle > b.groupTitle ? 1 : -1)
  }

  trackByGroup(_: number, group: IssueGroup) {
    return group.groupId
  }

  trackByIssue(index: number, issue: KanbanIssue) {
    return issue.id
  }

  loadMore() {
    if (this.issuePage && (this.issuePage!.page < this.issuePage!.pages)) {
      this.isPageMore$.next(true);
    }
  }

  onEnd(el: HTMLElement) {
    this.loadMore()
  }

  addHighlighted(id: number) {
    this.highlightedIds.add(id)
    setTimeout(() => this.zone.run(() => this.highlightedIds.delete(id)), this.highlightedTimeout)
  }
}
