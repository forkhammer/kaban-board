import { Component, DestroyRef, ElementRef, forwardRef, inject, Input, OnInit, ViewChild } from '@angular/core';
import { ControlValueAccessor, FormBuilder, FormGroup, NG_VALUE_ACCESSOR } from '@angular/forms';
import { BehaviorSubject, catchError, combineLatestWith, debounceTime, distinctUntilChanged, EMPTY, filter, map, of, pluck, switchMap } from 'rxjs';
import { Sprint, SprintStatus } from '../../models/sprint';
import { ToastService } from 'src/app/modules/core/services/toast.service';
import { SprintService } from '../../services/sprint.service';
import { faTimes, faRunning } from '@fortawesome/free-solid-svg-icons';
import { Pagination } from 'src/app/modules/core/models/base';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { HttpErrorResponse } from '@angular/common/module.d-CnjH8Dlt';
import { catchErrorMessages } from 'src/app/modules/core/tools/catch-error';
import {faClock} from '@fortawesome/free-regular-svg-icons';
import {faPlay, faStop} from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'app-sprint-select',
  templateUrl: './sprint-select.component.html',
  styleUrls: [
    '../../../ui/components/select/select.component.scss',
    '../../../ui/components/select/filter-form.scss',
    './sprint-select.component.scss',
  ],
  providers: [
    { provide: NG_VALUE_ACCESSOR, useExisting: forwardRef(() => SprintSelectComponent), multi: true },
  ],
  standalone: false,
})
export class SprintSelectComponent implements ControlValueAccessor, OnInit {
  protected fb = inject(FormBuilder)
  destroyRef = inject(DestroyRef)
  toast = inject(ToastService)
  sprintService = inject(SprintService)

  @Input() useSearch: boolean = false;
  @Input() useClear: boolean = false;
  @Input() title = '';
  @Input() all = true;
  @Input() nullTitle: string | null = null;
  @ViewChild('dropdown') dropdown!: ElementRef;

  SprintStatus = SprintStatus;
  faClock = faClock
  faPlay = faPlay
  faStop = faStop

  value = new BehaviorSubject<number | null>(null);
  valueModel$ = new BehaviorSubject<Sprint | null>(null);
  valuesModel: Sprint[] = [];
  valuesFilter = new BehaviorSubject<any>(null);
  private onChange: any;
  faTimes = faTimes;
  faRunning = faRunning;
  searchForm: FormGroup;
  protected errorValuesMessage: string | null = null;
  initialLoad$ = new BehaviorSubject<boolean>(false);

  get selectValue(): Sprint | null {
    return null;
  }

  set selectValue(item: Sprint | null) {
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
  set valueModel(value: Sprint | null) {
    this.valueModel$.next(value);
  }

  constructor() {
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
          const query = Object.assign({}, data, { search: searchControl ? searchControl.value : null, all: this.all });
          return this.sprintService.list(query)
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
          if (this.sprintService.usePagination) {
            return (data as Pagination<Sprint>).results;
          } else {
            return (data as Sprint[]);
          }
        }),
        takeUntilDestroyed(this.destroyRef)
      )
      .subscribe(data => {
        this.valuesModel = data;
      });

    this.searchForm.valueChanges
      .pipe(
        pluck('search'),
        debounceTime(500),
        switchMap(data => {
          this.errorValuesMessage = null;
          const query = Object.assign({}, this.valuesFilter.value, { search: data, all: this.all });
          return this.sprintService.list(query)
            .pipe(catchError((err: HttpErrorResponse) => {
              this.errorValuesMessage = err.statusText;
              return EMPTY;
            }));
        }),
        map((data: any) => {
          if (this.sprintService.usePagination) {
            return (data as Pagination<Sprint>).results;
          } else {
            return (data as Sprint[]);
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
              return this.sprintService.get(data);
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

  writeValue(value: number | null) {
    this.value.next(value);
    if (this.onChange) {
      this.onChange(value);
    }
  }

  registerOnChange(fn: any) {
    this.onChange = fn;
  }

  registerOnTouched(fn: any) {
  }

  getSearchText(): string | null {
    const searchControl = this.searchForm.get('search');
    return searchControl ? searchControl.value : null;
  }

  clearFilter(e: MouseEvent) {
    this.searchForm.patchValue({
      search: ''
    });
    return false;
  }

  select(e: MouseEvent, item: Sprint) {
    this.selectValue = item;
    (this.dropdown as any).close();
    return false;
  }

  clear(e: MouseEvent) {
    this.selectValue = null;
    (this.dropdown as any).close();
    e.stopPropagation();
    return false;
  }

  onOpenChange(open: boolean) {
    if (open) {
      this.initialLoad$.next(true);
    }
  }

  trackByItem(_: number, item: Sprint) {
    return item.id;
  }

  getStatusLabel(status: SprintStatus): string {
    switch (status) {
      case SprintStatus.RUNNING: return 'Активный';
      case SprintStatus.COMPLETED: return 'Завершён';
      case SprintStatus.WAITING: return 'Ожидает';
    }
  }
}
