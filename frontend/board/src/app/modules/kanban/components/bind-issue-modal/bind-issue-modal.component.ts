import { Component, DestroyRef, inject } from '@angular/core';
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
export class BindIssueModalComponent {
  modal = inject(NgbActiveModal);
  issueService = inject(IssueService);
  destroyRef = inject(DestroyRef);

  searchQuery = '';
  issues: KanbanIssue[] = [];
  isLoading = false;

  private search$ = new Subject<string>();

  constructor() {
    this.search$.pipe(
      debounceTime(300),
      distinctUntilChanged(),
      switchMap(query => {
        this.isLoading = true;
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

  onSearchChange(query: string) {
    this.search$.next(query);
  }

  select(issue: KanbanIssue) {
    this.modal.close(issue);
  }

  close() {
    this.modal.dismiss();
  }
}
