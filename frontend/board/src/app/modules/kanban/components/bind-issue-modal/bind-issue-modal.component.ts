import { Component, DestroyRef, ElementRef, inject, AfterViewInit, ViewChild, QueryList, ViewChildren } from '@angular/core';
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';
import { IssueService } from '../../services/issue.service';
import { KanbanIssue } from '../../models/kanban-issue';
import { Pagination } from 'src/app/modules/core/models/base';
import { Subject, debounceTime, distinctUntilChanged, switchMap, of, catchError } from 'rxjs';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Component({
  selector: 'app-bind-issue-modal',
  templateUrl: './bind-issue-modal.component.html',
  styleUrl: './bind-issue-modal.component.scss',
  standalone: false
})
export class BindIssueModalComponent implements AfterViewInit {
  @ViewChild('searchInput') searchInput!: ElementRef<HTMLInputElement>;
  @ViewChildren('issueItem') issueItems!: QueryList<ElementRef<HTMLElement>>;

  modal = inject(NgbActiveModal);
  issueService = inject(IssueService);
  destroyRef = inject(DestroyRef);

  searchQuery = '';
  issues: KanbanIssue[] = [];
  isLoading = false;
  activeIndex = -1;

  private search$ = new Subject<string>();

  constructor() {
    this.search$.pipe(
      debounceTime(300),
      distinctUntilChanged(),
      switchMap(query => {
        this.isLoading = true;
        this.activeIndex = -1;
        if (!query.trim()) {
          this.isLoading = false;
          return of(null);
        }
        return this.issueService.list({ search: query }).pipe(
          catchError(() => {
            this.isLoading = false;
            return of(null);
          })
        );
      }),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(result => {
      this.isLoading = false;
      if (result) {
        this.issues = (result as Pagination<KanbanIssue>).results;
      } else {
        this.issues = [];
      }
    });
  }

  ngAfterViewInit() {
    this.searchInput.nativeElement.focus();
  }

  onSearchChange(query: string) {
    this.search$.next(query);
  }

  onKeyDown(event: KeyboardEvent) {
    if (!this.issues.length) return;

    if (event.key === 'ArrowDown') {
      event.preventDefault();
      this.activeIndex = Math.min(this.activeIndex + 1, this.issues.length - 1);
      this.scrollActiveIntoView();
    } else if (event.key === 'ArrowUp') {
      event.preventDefault();
      this.activeIndex = Math.max(this.activeIndex - 1, 0);
      this.scrollActiveIntoView();
    } else if (event.key === 'Enter' && this.activeIndex >= 0) {
      event.preventDefault();
      this.select(this.issues[this.activeIndex]);
    }
  }

  private scrollActiveIntoView() {
    const items = this.issueItems.toArray();
    items[this.activeIndex]?.nativeElement.scrollIntoView({ block: 'nearest' });
  }

  select(issue: KanbanIssue) {
    this.modal.close(issue);
  }

  close() {
    this.modal.dismiss();
  }
}
