import { Component, OnInit, Input, ViewChild, ElementRef, forwardRef, OnDestroy } from '@angular/core';
import { faAngleDown, faAngleUp, faTimes } from '@fortawesome/free-solid-svg-icons';
import { NG_VALUE_ACCESSOR, ControlValueAccessor, FormGroup, FormBuilder } from '@angular/forms';
import {Subject, BehaviorSubject, of, EMPTY, distinctUntilChanged} from 'rxjs';
import { BaseService } from '../../../core/services/base.service';
import { BaseTitleModel, Pagination } from '../../../core/models/base';
import { takeUntil, switchMap, pluck, debounceTime, catchError, map } from 'rxjs/operators';
import { HttpErrorResponse } from '@angular/common/http';

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
  exportAs: 'selectModel'
})
export class SelectModelComponent implements ControlValueAccessor, OnInit, OnDestroy {

  @Input() useSearch: boolean = false;
  @Input() useClear: boolean = false;
  @Input() service!: BaseService<BaseTitleModel>;
  @Input() title = '';
  @Input() all = true;
  @Input() nullTitle: string | null = null;
  @Input() formatter: any = null;
  @Input() itemFormatter: any = null;
  @ViewChild('dropdown') dropdown!: ElementRef;

  value = new BehaviorSubject<string | number | null>(null);
  valueModel$ = new BehaviorSubject<BaseTitleModel | null>(null);
  valuesModel: BaseTitleModel[] = [];
  valuesFilter = new BehaviorSubject<any>(null);
  private onChange: any;
  faAngleDown = faAngleDown;
  faAngleUp = faAngleUp;
  faTimes = faTimes;
  searchForm: FormGroup;
  private destroy$ = new Subject();
  protected errorValuesMessage: string | null = null;

  get selectValue(): BaseTitleModel | null {
    return null;
  }

  set selectValue(item: BaseTitleModel | null) {
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

  constructor(
    protected fb: FormBuilder
  ) {
    this.searchForm = this.fb.group({
      search: ['']
    });
  }

  ngOnInit() {
    this.valuesFilter
      .pipe(
        switchMap(data => {
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
        takeUntil(this.destroy$),
      )
      .subscribe(data => {
        this.valuesModel = data;
      });

    // поиск
    this.searchForm.valueChanges
      .pipe(
          takeUntil(this.destroy$),
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
            return (data as Pagination<BaseTitleModel>).results;
          }),
      )
      .subscribe(data => {
        this.valuesModel = data;
      });

    this.value
      .pipe(
        distinctUntilChanged(),
        takeUntil(this.destroy$),
        switchMap(data => {
          if (data) {
            return this.service.get(data);
          } else {
            return [null];
          }
        })
      )
      .subscribe(data => {
        this.valueModel$.next(data);
      });
  }

  ngOnDestroy() {
    this.destroy$.next(null);
    this.destroy$.complete();
  }

  writeValue(value: string | number | null) {
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
    const searchCntrol = this.searchForm.get('search');
    return searchCntrol ? searchCntrol.value : null;
  }

  clearFilter(e: MouseEvent) {
    this.searchForm.patchValue({
      search: ''
    });
    return false;
  }

  select(e: MouseEvent, item: BaseTitleModel) {
    this.selectValue = item;
    (this.dropdown as any).close();
    return false;
  }

  clear(e: MouseEvent) {
    this.selectValue = null;
    (this.dropdown as any).close();
    return false;
  }

  getTitle(item: BaseTitleModel) {
    if (this.formatter) {
      return this.formatter(item);
    }
    return item.title;
  }

  getItemTitle(item: BaseTitleModel) {
    if (this.itemFormatter) {
      return this.itemFormatter(item);
    }
    return item.title;
  }

}
