import { Component, DestroyRef, ElementRef, forwardRef, inject, Input, OnInit, ViewChild } from '@angular/core';
import { ControlValueAccessor, FormBuilder, FormGroup, NG_VALUE_ACCESSOR } from '@angular/forms';
import { BehaviorSubject, catchError, combineLatestWith, debounceTime, distinctUntilChanged, EMPTY, filter, map, of, pluck, switchMap } from 'rxjs';
import { User } from '../../models/user';
import { ToastService } from 'src/app/modules/core/services/toast.service';
import { UserService } from '../../services/user.service';
import { faTimes } from '@fortawesome/free-solid-svg-icons';
import { faUser } from '@fortawesome/free-regular-svg-icons';
import { Pagination } from 'src/app/modules/core/models/base';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { HttpErrorResponse } from '@angular/common/module.d-CnjH8Dlt';
import { catchErrorMessages } from 'src/app/modules/core/tools/catch-error';

@Component({
  selector: 'app-assignee-select',
  templateUrl: './assignee-select.component.html',
  styleUrls: [
        '../../../ui/components/select/select.component.scss',
        '../../../ui/components/select/filter-form.scss',
        './assignee-select.component.scss',
    ],
  providers: [
    { provide: NG_VALUE_ACCESSOR, useExisting: forwardRef(() => AssigneeSelectComponent), multi: true },
  ],
  standalone: false,
})
export class AssigneeSelectComponent implements ControlValueAccessor, OnInit {
  protected fb = inject(FormBuilder)
  destroyRef = inject(DestroyRef)
  toast = inject(ToastService)
  userService = inject(UserService)

  @Input() useSearch: boolean = false;
  @Input() useClear: boolean = false;
  @Input() title = '';
  @Input() all = true;
  @Input() nullTitle: string | null = null;
  @Input() formatter: any = null;
  @Input() itemFormatter: any = null;
  @ViewChild('dropdown') dropdown!: ElementRef;
  @ViewChild('searchInput') searchInput!: ElementRef<HTMLInputElement>;

  value = new BehaviorSubject<number | null>(null);
  valueModel$ = new BehaviorSubject<User | null>(null);
  valuesModel: User[] = [];
  valuesFilter = new BehaviorSubject<any>(null);
  private onChange: any;
  faTimes = faTimes;
  faUser = faUser;
  searchForm: FormGroup;
  protected errorValuesMessage: string | null = null;
  initialLoad$ = new BehaviorSubject<boolean>(false);

  get selectValue(): User | null {
    return null;
  }

  set selectValue(item: User | null) {
    if (item) {
      this.writeValue(item.id);
    } else {
      this.writeValue(null);
    }
  }

  @Input()
  set filter(value: any) {
    this.valuesFilter.next(value);
  }

  @Input()
  set valueModel(value: User | null) {
    this.valueModel$.next(value);
  }

  constructor(
  ) {
    this.searchForm = this.fb.group({
      search: ['']
    });

  }

  ngOnInit() {
    this.valuesFilter
      .pipe(
        combineLatestWith(this.initialLoad$.pipe(filter(Boolean), distinctUntilChanged())),
        switchMap(([data, _]) => {
          this.errorValuesMessage = null;
          const searchControl = this.searchForm.get('search');
          const query = Object.assign({}, data, {search: searchControl ? searchControl.value : null, all: this.all});
          return this.userService.list(query)
            .pipe(catchError((err: HttpErrorResponse) => {
              this.errorValuesMessage = err.statusText;
              if (this.all) {
                return of([]);
              }
              return of([{
                page: 1,
                results: []
              }]);
            }));
        }),
        map((data: any) => {
          if (this.userService.usePagination) {
            return (data as Pagination<User>).results;
          } else {
            return (data as User[]);
          }
        }),
        takeUntilDestroyed(this.destroyRef)
      )
      .subscribe(data => {
        this.valuesModel = data;
      });

    // поиск
    this.searchForm.valueChanges
      .pipe(
        pluck('search'),
        debounceTime(500),
        switchMap(data => {
          this.errorValuesMessage = null;
          const query = Object.assign({}, this.valuesFilter.value, {search: data, all: this.all});
          return this.userService.list(query)
          .pipe(catchError((err: HttpErrorResponse) => {
            this.errorValuesMessage = err.statusText;
            return EMPTY;
          }));
        }),
        map((data: any) => {
          if (this.userService.usePagination) {
            return (data as Pagination<User>).results;
          } else {
            return (data as User[]);
          }
        }),
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe(data => {
        this.valuesModel = data;
      });

    this.value
      .pipe(
        distinctUntilChanged(),
        debounceTime(10),
        switchMap(data => {
          if (data) {
            if (this.valueModel$.value?.id === data) {
              return of(this.valueModel$.value);
            } else {
              return this.userService.get(data);
            }
          } else {
            return [null];
          }
        }),
        catchErrorMessages(this.toast),
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe(data => {
        this.valueModel$.next(data);
      });
  }

  private emitChange(value: number | null): void {
    if (this.onChange) {
      this.onChange(value);
    }
  }

  writeValue(value: number | null) {
    this.value.next(value);
  }

  registerOnChange(fn: any) {
    this.onChange = fn;
  }

  registerOnTouched(fn: any) {

  }

  getSearchText(): string | null {
    const searchCntrol = this.searchForm.get('search');
    return searchCntrol ? searchCntrol.value : null;
  }

  clearFilter(e: MouseEvent) {
    this.searchForm.patchValue({
      search: ''
    });
    return false;
  }

  select(e: MouseEvent, item: User) {
    this.selectValue = item;
    this.emitChange(this.value.value);
    (this.dropdown as any).close();
    return false;
  }

  clear(e: MouseEvent) {
    this.selectValue = null;
    this.emitChange(null);
    (this.dropdown as any).close();
    return false;
  }

  onOpenChange(open: boolean) {
    if (open) {
      this.initialLoad$.next(true);
      if (this.useSearch && this.searchInput) {
        setTimeout(() => this.searchInput.nativeElement.focus(), 0);
      }
    }
  }

  trackByItem(_: number, item: User) {
    return item.id;
  }
}
