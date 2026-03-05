import { Component, DestroyRef, ElementRef, inject, AfterViewInit, ViewChild, QueryList, ViewChildren } from '@angular/core';
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';
import { IssueService } from '../../services/issue.service';
import { ProjectService } from '../../services/project.service';
import { KanbanIssue } from '../../models/kanban-issue';
import { Pagination } from 'src/app/modules/core/models/base';
import { Subject, debounceTime, switchMap, of, catchError } from 'rxjs';
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
  projectService = inject(ProjectService);
  destroyRef = inject(DestroyRef);

  searchQuery = '';
  selectedProjectId: number | null = null;
  issues: KanbanIssue[] = [];
  isLoading = false;
  activeIndex = -1;

  private search$ = new Subject<void>();

  constructor() {
    this.search$.pipe(
      debounceTime(300),
      switchMap(() => {
        this.isLoading = true;
        this.activeIndex = -1;
        if (!this.searchQuery.trim()) {
          this.isLoading = false;
          return of(null);
        }
        const query: any = { search: this.searchQuery };
        if (this.selectedProjectId) {
          query['project'] = this.selectedProjectId;
        }
        return this.issueService.list(query).pipe(
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

  onSearchChange() {
    this.search$.next();
  }

  onProjectChange() {
    this.search$.next();
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
