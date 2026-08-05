import {
  Component,
  DestroyRef,
  EventEmitter,
  inject,
  Input,
  NgZone,
  Output,
} from "@angular/core";
import {
  BehaviorSubject,
  catchError,
  combineLatestWith,
  debounceTime,
  distinctUntilChanged,
  EMPTY,
  filter,
  merge,
  of,
  Subject,
  switchMap,
  timer,
} from "rxjs";
import { KanbanUser } from "../../models/kanban-user";
import { Team } from "../../models/team";
import { Sprint } from "../../models/sprint";
import { environment } from "src/environments/environment";
import { isEqual } from "lodash";
import { catchErrorMessages } from "src/app/modules/core/tools/catch-error";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import {
  BIND_STATUS_WEIGHT,
  IssueGroup,
  KanbanIssue,
  PRIORITY_WEIGHT,
} from "../../models/kanban-issue";
import { IssueService } from "../../services/issue.service";
import { ToastService } from "src/app/modules/core/services/toast.service";
import {
  faPlus,
  faSort,
  faSortDown,
  faSortUp,
} from "@fortawesome/free-solid-svg-icons";
import { Pagination } from "src/app/modules/core/models/base";
import { AccountService } from "src/app/modules/core/services/account.service";
import { IssueBindingService } from "../../services/issue-binding.service";
import { BindIssueModalService } from "../../services/bind-issue-modal.service";
import { CdkDragDrop, moveItemInArray } from "@angular/cdk/drag-drop";
import { IssueBindingModalService } from "../../services/issue-binding-modal.service";
import {
  BindingDeletedEventData,
  BindingUpdatedEventData,
  SprintWebsocketService,
} from "../../services/sprint-websocket.service";
import { Group } from "../../models/group";

const SORTABLE_COLUMNS = [
  "task",
  "assignee",
  "epic",
  "release",
  "priority",
  "status",
] as const;
type SortableColumn = (typeof SORTABLE_COLUMNS)[number];

@Component({
  selector: "app-issue-table",
  standalone: false,
  templateUrl: "./issue-table.component.html",
  styleUrl: "./issue-table.component.scss",
})
export class IssueTableComponent {
  private issueService = inject(IssueService);
  private issueBindingService = inject(IssueBindingService);
  private toast = inject(ToastService);
  public accountService = inject(AccountService);
  private bindIssueModal = inject(BindIssueModalService);
  private destroyRef = inject(DestroyRef);
  private issueBindingModal = inject(IssueBindingModalService);
  private sprintWs = inject(SprintWebsocketService);
  private zone = inject(NgZone);

  faPlus = faPlus;

  @Output() sprintStateChanged = new EventEmitter<void>();

  public user$ = new BehaviorSubject<KanbanUser | null | undefined>(null);
  public team$ = new BehaviorSubject<Team | null | undefined>(null);
  public sprint$ = new BehaviorSubject<Sprint | null | undefined>(null);
  public group$ = new BehaviorSubject<Group | null | undefined>(null);
  private timer$ = timer(0, environment.autoUpdateIssuesMin * 60 * 1000);
  private wsReload$ = new BehaviorSubject<number>(0);
  public issues: KanbanIssue[] = [];
  public groupedIssues: IssueGroup[] = [];
  public issuePage: Pagination<KanbanIssue> | null = null;
  public isLoading = false;
  public search$ = new BehaviorSubject<string | null>(null);
  isPageMore$ = new BehaviorSubject<boolean>(false);
  isLoadMore$ = new BehaviorSubject<boolean>(false);
  highlightedIds = new Set<number>();
  highlightedTimeout = 1000;

  faSort = faSort;
  faSortUp = faSortUp;
  faSortDown = faSortDown;

  readonly SORTABLE_COLUMNS = SORTABLE_COLUMNS;

  private static readonly SORT_STORAGE_KEY = 'issue-table-sort';

  sortColumn: SortableColumn | null = this.loadSortColumn();
  sortDirection: "asc" | "desc" = this.loadSortDirection();

  @Input() set user(value: KanbanUser | undefined | null) {
    this.user$.next(value);
  }

  @Input() set team(value: Team | undefined | null) {
    this.team$.next(value);
  }

  @Input() set sprint(value: Sprint | undefined | null) {
    this.sprint$.next(value);
  }

  @Input() set group(value: Group | undefined | null) {
    this.group$.next(value);
  }

  @Input() set search(value: string | undefined | null) {
    this.search$.next(value ?? null);
  }

  constructor() {
    this.isPageMore$
      .pipe(filter(Boolean), takeUntilDestroyed())
      .subscribe((data) => {
        this.isLoadMore$.next(data);
      });

    this.sprint$
      .pipe(
        distinctUntilChanged((a, b) => a?.id === b?.id),
        takeUntilDestroyed(),
      )
      .subscribe((sprint) => {
        if (sprint) {
          this.sprintWs.connect(sprint.id);
        } else {
          this.sprintWs.disconnect();
        }
      });

    this.sprintWs.events$.pipe(takeUntilDestroyed()).subscribe((event) => {
      switch (event.type) {
        case "binding_updated":
          this.onBindingUpdated(event.data as BindingUpdatedEventData);
          break;
        case "binding_deleted":
          this.onBindingDeleted(event.data as BindingDeletedEventData);
          break;
        case "binding_created":
          this.onBindingCreated(event.data as BindingUpdatedEventData);
          break;
        case "binding_ordering":
          this.reloadIssues();
          break;
      }
    });

    this.destroyRef.onDestroy(() => this.sprintWs.disconnect());

    merge(this.timer$, this.wsReload$)
      .pipe(
        combineLatestWith(
          this.user$,
          this.team$,
          this.sprint$,
          this.group$,
          this.search$,
        ),
        debounceTime(100),
        filter(([_, user, team, sprint, group, search]) => {
          return !!team || !!search;
        }),
        distinctUntilChanged(isEqual),
        debounceTime(1),
        combineLatestWith(this.isLoadMore$),
        switchMap(([data, isLoadingMore]) => {
          const [_, user, team, sprint, group, search] = data as [
            any,
            KanbanUser,
            Team,
            Sprint,
            Group,
            string,
          ];
          const query: Record<string, any> = {
            limit: sprint ? 10000 : 50,
            page: this.isPageMore$.value ? this.issuePage!.page + 1 : 1,
          };
          if (team) {
            query["team"] = team.id;
          }
          if (user) {
            query["assignee"] = user.id;
          }
          if (sprint) {
            query["sprint"] = sprint.id;
          }
          if (group) {
            query["group"] = group.id;
          }
          if (search) {
            query["search"] = search;
          }
          this.isLoading = true;

          if (sprint) {
            return this.issueBindingService.list(query).pipe(
              combineLatestWith(of(this.isPageMore$.value)),
              catchErrorMessages(this.toast, () => (this.isLoading = false)),
            );
          } else {
            return this.issueService.list(query).pipe(
              combineLatestWith(of(this.isPageMore$.value)),
              catchErrorMessages(this.toast, () => (this.isLoading = false)),
            );
          }
        }),
        takeUntilDestroyed(),
      )
      .subscribe(([data, isLoadingMore]) => {
        this.isLoading = false;

        if (isLoadingMore && this.issuePage) {
          this.issuePage = data as Pagination<KanbanIssue>;
          this.setIssues(this.issues.concat(this.issuePage.results));
          this.isPageMore$.next(false);
        } else {
          this.issuePage = data as Pagination<KanbanIssue>;
          this.setIssues(
            [...this.issuePage.results].sort((a, b) => {
              const oa = a.order ?? "";
              const ob = b.order ?? "";
              if (oa === "" && ob === "") return 0;
              if (oa === "") return 1;
              if (ob === "") return -1;
              return oa.localeCompare(ob);
            }),
          );
        }
      });
  }

  private onBindingUpdated(data: BindingUpdatedEventData): void {
    this.sprintStateChanged.emit();

    if (data.account_id === this.accountService.user$.value?.id) return;

    const idx = this.issues.findIndex((i) => i.bindingId === data.id);
    if (idx === -1) return;

    this.issueBindingService
      .get(data.id)
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe((updated) => {
        Object.assign(this.issues[idx], updated);
        this.groupedIssues = this.getGroupedIssues();
        this.addHighlighted(data.id);
      });
  }

  private onBindingDeleted(data: BindingDeletedEventData): void {
    this.setIssues(this.issues.filter((i) => i.bindingId !== data.id));
    this.sprintStateChanged.emit();
    this.groupedIssues = this.getGroupedIssues();
  }

  private reloadIssues(): void {
    this.wsReload$.next(this.wsReload$.value + 1);
  }

  private onBindingCreated(data: BindingUpdatedEventData): void {
    this.wsReload$.next(this.wsReload$.value + 1);
    this.sprintStateChanged.emit();
    this.groupedIssues = this.getGroupedIssues();
    if (data.account_id !== this.accountService.user$.value?.id) {
      this.addHighlighted(data.id);
    }
  }

  appendIssue() {
    this.bindIssueModal.show().then(
      (result) => {
        if (!result || !this.sprint$.value) {
          return;
        }
        const issues = Array.isArray(result) ? result : [result];
        const ids = issues.map((i) => +i.id);
        this.issueService
          .bindToSprint(
            ids,
            this.sprint$.value.id,
            this.user$.value?.id ?? null,
          )
          .pipe(
            catchErrorMessages(this.toast),
            takeUntilDestroyed(this.destroyRef),
          )
          .subscribe((data) => {
            for (const error of data.errors) {
              this.toast.show(this.toast.createErrorToast(error));
            }
            for (const issue of data.results) {
              this.issues.push(issue);
            }
            this.groupedIssues = this.getGroupedIssues();
            this.sprintStateChanged.emit();
          });
      },
      () => {},
    );
  }

  createIssue() {
    this.issueBindingModal
      .show({
        assigneeId: this.user$.value?.id ?? null,
        sprintId: this.sprint$.value?.id ?? null,
      })
      .then(
        (result) => {
          if (result) {
            this.issues.push(result);
            this.groupedIssues = this.getGroupedIssues();
            this.sprintStateChanged.emit();
          }
        },
        (err) => {},
      );
  }

  unbindIssue(bindingId: number) {
    this.setIssues(
      this.issues.filter((issue) => issue.bindingId !== bindingId),
    );
    this.sprintStateChanged.emit();
  }

  toggleSort(column: SortableColumn): void {
    if (this.sortColumn === column) {
      if (this.sortDirection === "asc") {
        this.sortDirection = "desc";
      } else {
        this.sortColumn = "assignee";
        this.sortDirection = "asc";
      }
    } else {
      this.sortColumn = column;
      this.sortDirection = "asc";
    }
    this.saveSortState();
  }

  onSaveIssue(issue: KanbanIssue) {
    const idx = this.issues.findIndex((i) => i.bindingId === issue.bindingId);
    const originalIssue = idx !== -1 ? { ...this.issues[idx] } : null;

    this.issueBindingService
      .save(issue)
      .pipe(
        catchErrorMessages(this.toast, () => {
          if (idx !== -1 && originalIssue) {
            this.issues[idx] = originalIssue;
            this.groupedIssues = this.getGroupedIssues();
          }
        }),
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe((data) => {
        if (idx !== -1) {
          this.issues[idx] = { ...this.issues[idx], ...data };
          this.groupedIssues = this.getGroupedIssues();
        }
      });
  }

  onDeleteIssue(issue: KanbanIssue) {
    this.issueBindingService
      .delete(issue)
      .pipe(catchErrorMessages(this.toast), takeUntilDestroyed(this.destroyRef))
      .subscribe(() => {
        this.unbindIssue(issue.bindingId!);
      });
  }

  dropIssue(event: CdkDragDrop<KanbanIssue[]>): void {
    const flat = this.displayedIssueGroups.flatMap((g) => g.issues);
    moveItemInArray(flat, event.previousIndex, event.currentIndex);
    this.recalculateOrder(flat);
    this.setIssues(flat);
    this.issueBindingService
      .saveOrdering(flat)
      .pipe(catchErrorMessages(this.toast), takeUntilDestroyed(this.destroyRef))
      .subscribe();
  }

  private recalculateOrder(issues: KanbanIssue[]): void {
    issues.forEach((issue, index) => {
      issue.order = String(index).padStart(8, "0");
    });
  }

  private loadSortColumn(): SortableColumn | null {
    try {
      const raw = localStorage.getItem(IssueTableComponent.SORT_STORAGE_KEY);
      if (raw) {
        const parsed = JSON.parse(raw);
        if (parsed.column && SORTABLE_COLUMNS.includes(parsed.column)) {
          return parsed.column as SortableColumn;
        }
      }
    } catch { /* localStorage not available or corrupt */ }
    return 'assignee';
  }

  private loadSortDirection(): "asc" | "desc" {
    try {
      const raw = localStorage.getItem(IssueTableComponent.SORT_STORAGE_KEY);
      if (raw) {
        const parsed = JSON.parse(raw);
        if (parsed.direction === 'asc' || parsed.direction === 'desc') {
          return parsed.direction;
        }
      }
    } catch { /* localStorage not available or corrupt */ }
    return 'asc';
  }

  private saveSortState(): void {
    try {
      localStorage.setItem(
        IssueTableComponent.SORT_STORAGE_KEY,
        JSON.stringify({ column: this.sortColumn, direction: this.sortDirection }),
      );
    } catch { /* localStorage not available */ }
  }

  private sortIssues(issues: KanbanIssue[]): KanbanIssue[] {
    if (!this.sortColumn) return issues;

    const sorted = [...issues].sort((a, b) => {
      const cmp = this.compareByColumn(a, b, this.sortColumn!);
      if (cmp !== 0) return this.sortDirection === "asc" ? cmp : -cmp;

      const oa = a.order ?? "";
      const ob = b.order ?? "";
      if (oa === "" && ob === "") return 0;
      if (oa === "") return 1;
      if (ob === "") return -1;
      return oa.localeCompare(ob);
    });

    return sorted;
  }

  private compareByColumn(
    a: KanbanIssue,
    b: KanbanIssue,
    column: SortableColumn,
  ): number {
    switch (column) {
      case "task":
        return 0;
      case "assignee": {
        const an = a.assignee?.name?.toLowerCase() ?? "";
        const bn = b.assignee?.name?.toLowerCase() ?? "";
        return an.localeCompare(bn);
      }
      case "epic": {
        const ae = a.epic?.title?.toLowerCase() ?? "";
        const be = b.epic?.title?.toLowerCase() ?? "";
        return ae.localeCompare(be);
      }
      case "release": {
        const ar = a.release?.title?.toLowerCase() ?? "";
        const br = b.release?.title?.toLowerCase() ?? "";
        return ar.localeCompare(br);
      }
      case "priority": {
        const pw = PRIORITY_WEIGHT;
        const pa = a.priority ? (pw[a.priority] ?? 999) : 999;
        const pb = b.priority ? (pw[b.priority] ?? 999) : 999;
        return pa - pb;
      }
      case "status": {
        const sw = BIND_STATUS_WEIGHT;
        const sa = a.bindStatus ? (sw[a.bindStatus] ?? 999) : 999;
        const sb = b.bindStatus ? (sw[b.bindStatus] ?? 999) : 999;
        return sa - sb;
      }
      default:
        return 0;
    }
  }


  setIssues(issues: KanbanIssue[]): void {
    this.issues = issues;
    this.groupedIssues = this.getGroupedIssues();
  }

  getGroupedIssues(): IssueGroup[] {
    const map = new Map<number, IssueGroup>();
    const NO_GROUP_ID = -1;

    const teamGroupIds = new Set(
      (this.team$.value?.groups ?? []).map((g) => g.id),
    );

    for (const issue of this.issues) {
      const group = issue.assignee?.groups?.find((g) => teamGroupIds.has(g.id));
      const key = group?.id ?? NO_GROUP_ID;
      if (!map.has(key)) {
        map.set(key, {
          groupId: key,
          groupTitle: group?.title ?? "Без группы",
          issues: [],
        });
      }
      map.get(key)!.issues.push(issue);
    }

    return Array.from(map.values()).sort((a, b) =>
      a.groupTitle > b.groupTitle ? 1 : -1,
    );
  }

  get displayedIssueGroups(): IssueGroup[] {
    if (this.sortColumn === null) return this.groupedIssues;
    return this.groupedIssues.map((g) => ({
      ...g,
      issues: this.sortIssues(g.issues),
    }));
  }

  trackByGroup(_: number, group: IssueGroup) {
    return group.groupId;
  }

  trackByIssue(index: number, issue: KanbanIssue) {
    return issue.bindingId;
  }

  loadMore() {
    if (this.issuePage && this.issuePage!.page < this.issuePage!.pages) {
      this.isPageMore$.next(true);
    }
  }

  onEnd(el: HTMLElement) {
    this.loadMore();
  }

  addHighlighted(id: number) {
    this.highlightedIds.add(id);
    setTimeout(
      () => this.zone.run(() => this.highlightedIds.delete(id)),
      this.highlightedTimeout,
    );
  }
}
