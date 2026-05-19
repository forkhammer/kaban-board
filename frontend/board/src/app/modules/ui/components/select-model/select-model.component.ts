import { Component, Input, ViewChild, ElementRef, forwardRef, inject, OnInit, DestroyRef } from '@angular/core';
import { faTimes } from '@fortawesome/free-solid-svg-icons';
import { NG_VALUE_ACCESSOR, ControlValueAccessor, FormGroup, FormBuilder } from '@angular/forms';
import { BehaviorSubject, of, EMPTY, distinctUntilChanged} from 'rxjs';
import { BaseService } from '../../../core/services/base.service';
import { BaseModel, BaseTitleModel, Pagination } from '../../../core/models/base';
import { switchMap, pluck, debounceTime, catchError, map, combineLatestWith, filter } from 'rxjs/operators';
import { HttpErrorResponse } from '@angular/common/http';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { catchErrorMessages } from 'src/app/modules/core/tools/catch-error';
import { ToastService } from 'src/app/modules/core/services/toast.service';

@Component({
    selector: 'app-select-model',
    templateUrl: './select-model.component.html',
    styleUrls: [
        './select-model.component.scss',
        '../select/select.component.scss',
        '../select/filter-form.scss'
    ],
    providers: [
        { provide: NG_VALUE_ACCESSOR, useExisting: forwardRef(() => SelectModelComponent), multi: true },
    ],
    exportAs: 'selectModel',
    standalone: false
})
export class SelectModelComponent implements ControlValueAccessor, OnInit {
  protected fb = inject(FormBuilder)
  destroyRef = inject(DestroyRef)
  toast = inject(ToastService)

  @Input() useSearch: boolean = false;
  @Input() useClear: boolean = false;
  @Input() service!: BaseService<BaseModel>;
  @Input() title = '';
  @Input() all = true;
  @Input() nullTitle: string | null = null;
  @Input() formatter: any = null;
  @Input() itemFormatter: any = null;
  @ViewChild('dropdown') dropdown!: ElementRef;
  @ViewChild('searchInput') searchInput!: ElementRef<HTMLInputElement>;

  value = new BehaviorSubject<string | number | null>(null);
  valueModel$ = new BehaviorSubject<BaseModel | null>(null);
  valuesModel: BaseModel[] = [];
  valuesFilter = new BehaviorSubject<any>(null);
  private onChange: any;
  faTimes = faTimes;
  searchForm: FormGroup;
  protected errorValuesMessage: string | null = null;
  initialLoad$ = new BehaviorSubject<boolean>(false);
  reload$ = new BehaviorSubject<null>(null);

  get selectValue(): BaseModel | null {
    return null;
  }

  set selectValue(item: BaseModel | null) {
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
  set valueModel(value: BaseModel | null) {
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
        combineLatestWith(this.initialLoad$.pipe(filter(Boolean), distinctUntilChanged()), this.reload$),
        switchMap(([data, _, _1]) => {
          this.errorValuesMessage = null;
          const searchControl = this.searchForm.get('search');
          const query = Object.assign({}, data, {search: searchControl ? searchControl.value : null, all: this.all});
          return this.service.list(query)
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
          if (this.service.usePagination) {
            return (data as Pagination<BaseTitleModel>).results;
          } else {
            return (data as BaseTitleModel[]);
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
          return this.service.list(query)
          .pipe(catchError((err: HttpErrorResponse) => {
            this.errorValuesMessage = err.statusText;
            return EMPTY;
          }));
        }),
        map((data: any) => {
          if (this.service.usePagination) {
            return (data as Pagination<BaseModel>).results;
          } else {
            return (data as BaseModel[]);
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
              return this.service.get(data);
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

  private emitChange(value: string | number | null): void {
    if (this.onChange) {
      this.onChange(value);
    }
  }

  writeValue(value: string | number | null) {
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

  select(e: MouseEvent, item: BaseModel) {
    this.selectValue = item;
    this.emitChange(this.value.value);
    (this.dropdown as any).close();
    return false;
  }

  clear(e: MouseEvent) {
    this.selectValue = null;
    this.emitChange(null);
    (this.dropdown as any).close();
    e.preventDefault();
    return false;
  }

  getTitle(item: BaseModel) {
    if (this.formatter) {
      return this.formatter(item);
    }
    return (item as BaseTitleModel).title;
  }

  getItemTitle(item: BaseModel) {
    if (this.itemFormatter) {
      return this.itemFormatter(item);
    }
    return (item as BaseTitleModel).title;
  }

  onOpenChange(open: boolean) {
    if (open) {
      this.initialLoad$.next(true);
      if (this.useSearch && this.searchInput) {
        setTimeout(() => this.searchInput.nativeElement.focus(), 0);
      }
    }
  }

  trackByItem(_: number, item: BaseModel) {
    return item.id;
  }

  reload() {
    this.reload$.next(null)
  }

}
