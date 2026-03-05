import { Component, DestroyRef, inject, Input } from '@angular/core';
import { BehaviorSubject, combineLatestWith, debounceTime, distinctUntilChanged, filter, of, switchMap, tap, timer } from 'rxjs';
import { KanbanUser } from '../../models/kanban-user';
import { Team } from '../../models/team';
import { Sprint } from '../../models/sprint';
import { environment } from 'src/environments/environment';
import { isEqual } from 'lodash';
import { catchErrorMessages } from 'src/app/modules/core/tools/catch-error';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { KanbanIssue } from '../../models/kanban-issue';
import { IssueService } from '../../services/issue.service';
import { ToastService } from 'src/app/modules/core/services/toast.service';
import { faPlus } from '@fortawesome/free-solid-svg-icons'
import { Pagination } from 'src/app/modules/core/models/base';
import { AccountService } from 'src/app/modules/core/services/account.service';
import { IssueBindingService } from '../../services/issue-binding.service';
import { BindIssueModalService } from '../../services/bind-issue-modal.service';

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

  faPlus = faPlus

  public user$ = new BehaviorSubject<KanbanUser | null | undefined>(null)
  public team$ = new BehaviorSubject<Team | null | undefined>(null)
  public sprint$ = new BehaviorSubject<Sprint | null | undefined>(null)
  private timer$ = timer(0, environment.autoUpdateIssuesMin * 60 * 1000)
  public issues: KanbanIssue[] = []
  public issuePage: Pagination<KanbanIssue> | null = null
  public isLoading = false
  isPageMore$ = new BehaviorSubject<boolean>(false);
  isLoadMore$ = new BehaviorSubject<boolean>(false);

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

    this.timer$.pipe(
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
        this.issues = this.issues.concat(this.issuePage.results)
        this.isPageMore$.next(false)
      } else {
        this.issuePage = (data as Pagination<KanbanIssue>)
        this.issues = this.issuePage.results
      }
    })
  }

  appendIssue() {
    this.bindIssueModal.show().then(issue => {
      if (issue && this.sprint$.value) {
        this.issueService.bindToSprint((issue as KanbanIssue).id, this.sprint$.value.id, this.user$.value?.id ?? null).pipe(
          catchErrorMessages(this.toast),
          takeUntilDestroyed(this.destroyRef),
        ).subscribe(data => {
          this.issues.push(data)
        })
      }
    }, () => {})
  }

  unbindIssue(bindingId: number) {
    this.issues = this.issues.filter(issue => issue.bindingId !== bindingId)
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
}
