import { Component, OnInit, Input, forwardRef, OnDestroy } from '@angular/core';
import { faAngleDown, faAngleUp, faTimes, faCheck } from '@fortawesome/free-solid-svg-icons';
import { NG_VALUE_ACCESSOR, ControlValueAccessor, FormGroup, FormBuilder } from '@angular/forms';
import { Subject, BehaviorSubject, of, EMPTY } from 'rxjs';
import { takeUntil, switchMap, pluck, debounceTime, catchError, map } from 'rxjs/operators';
import { HttpErrorResponse } from '@angular/common/http';
import { BaseService } from "../../../core/services/base.service";
import { BaseTitleModel, Pagination } from "../../../core/models/base";

@Component({
  selector: 'app-select-model-multiple',
  templateUrl: './select-model-multiple.component.html',
  styleUrls: [
    '../select-multiple/select-multiple.component.scss',
    './select-model-multiple.component.scss',
  ],
  providers: [
    { provide: NG_VALUE_ACCESSOR, useExisting: forwardRef(() => SelectModelMultipleComponent), multi: true },
  ],
  exportAs: 'selectModel'
})
export class SelectModelMultipleComponent implements ControlValueAccessor, OnInit, OnDestroy {

  faAngleDown = faAngleDown;
  faAngleUp = faAngleUp;
  faTimes = faTimes;
  faCheck = faCheck;

  @Input() useSearch: boolean = false;
  @Input() useClear: boolean = false;
  @Input() service!: BaseService<BaseTitleModel>;
  @Input() title = '';
  @Input() all = true;
  @Input() nullTitle: string | null = null;
  @Input() formatter: any = null;

  value = new BehaviorSubject<(number | string)[]>([]);
  valuesModel: BaseTitleModel[] = [];
  valuesFilter = new BehaviorSubject<any>(null);
  private onChange: any;
  searchForm: FormGroup;
  private destroy$ = new Subject();
  protected errorValuesMessage: string | null = null;

  get selectValue(): (number | string)[] {
    return this.value.value;
  }

  set selectValue(item: (number | string)[]) {
    if (item) {
      this.writeValue(item);
    } else {
      this.writeValue([]);
    }
  }

  @Input()
  set filter(value: any) {
    this.valuesFilter.next(value);
  }

  constructor(
    private fb: FormBuilder
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
            if (this.service.usePagination) {
              return (data as Pagination<BaseTitleModel>).results;
            } else {
              return (data as BaseTitleModel[]);
            }
          }),
      )
      .subscribe(data => {
        this.valuesModel = data;
      });
  }

  ngOnDestroy() {
    this.destroy$.next(null);
    this.destroy$.complete();
  }

  writeValue(value: (number | string)[]) {
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

  getSearchText(): string {
    const searchControl = this.searchForm.get('search')
    return searchControl ? searchControl.value : null;
  }

  clearFilter(e: MouseEvent) {
    this.searchForm.patchValue({
      search: ''
    });
    return false;
  }

  select(e: MouseEvent, item: BaseTitleModel) {
    if (this.selectValue.indexOf(item.id) === -1) {
      this.selectValue = this.selectValue.concat(item.id);
    } else {
      this.selectValue = this.selectValue.filter(value => value !== item.id);
    }
    e.preventDefault();
    return false;
  }

  clear(e: MouseEvent) {
    this.selectValue = [];
    return false;
  }

  deselect(e: MouseEvent, item: BaseTitleModel) {
    if (this.selectValue.indexOf(item.id) > -1) {
      const value = [...this.selectValue];
      value.splice(value.indexOf(item.id), 1);
      this.selectValue = value;
    }
    e.preventDefault();
    return false;
  }

  getTitle(item: BaseTitleModel) {
    if (this.formatter) {
      return this.formatter(item);
    }
    return item.title;
  }

}
