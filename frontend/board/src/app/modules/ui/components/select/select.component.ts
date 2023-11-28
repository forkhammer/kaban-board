import { Component, forwardRef, Input, OnDestroy, ViewChild, ElementRef } from '@angular/core';
import { ControlValueAccessor, NG_VALUE_ACCESSOR, FormBuilder, FormGroup } from '@angular/forms';
import { faAngleDown, faAngleUp, faTimes } from '@fortawesome/free-solid-svg-icons';
import { Subject } from 'rxjs';
import { SelectValue, SelectValueIdentity } from '../../models/select-value';


@Component({
  selector: 'app-select',
  templateUrl: './select.component.html',
  styleUrls: ['./select.component.scss', './filter-form.scss'],
  providers: [
    { provide: NG_VALUE_ACCESSOR, useExisting: forwardRef(() => SelectComponent), multi: true },
  ]
})
export class SelectComponent implements ControlValueAccessor, OnDestroy {

  @Input() values: SelectValue[] = [];
  @Input() title = '';
  @Input() useSearch: boolean = false;
  @Input() useClear: boolean = false;
  @Input() nullTitle: string | null = null;
  @Input() formatter: any = null;
  @Input() itemFormatter: any = null;
  @ViewChild('dropdown') dropdown!: ElementRef;

  public value: SelectValueIdentity | null = null;
  private onChange: any;
  faAngleDown = faAngleDown;
  faAngleUp = faAngleUp;
  faTimes = faTimes;
  searchForm: FormGroup;
  private destroy$ = new Subject();

  get selectValue(): SelectValue | null {
    const result = this.values.filter(item => item.id === this.value);
    return result.length > 0 ? result[0] : null;
  }

  set selectValue(item: SelectValue | null) {
    if (item) {
      this.writeValue(item.id);
    } else {
      this.writeValue(null);
    }
  }

  constructor(
    private fb: FormBuilder
  ) {
    this.searchForm = this.fb.group({
      search: ['']
    });
  }

  ngOnDestroy() {
    this.destroy$.next(null);
    this.destroy$.complete();
  }

  writeValue(value: SelectValueIdentity | null) {
    this.value = value;
    if (this.onChange) {
      this.onChange(this.value);
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

  select(e: MouseEvent, item: SelectValue) {
    this.selectValue = item;
    if (this.dropdown) {
      (this.dropdown as any).close();
    }
    return false;
  }

  clear(e: MouseEvent) {
    this.selectValue = null;
    if (this.dropdown) {
      (this.dropdown as any).close();
    }
    return false;
  }

  getTitle(item: SelectValue) {
    if (this.formatter) {
      return this.formatter(item);
    }
    return item.title;
  }

  getItemTitle(item: SelectValue) {
    if (this.itemFormatter) {
      return this.itemFormatter(item);
    }
    return item.title;
  }
}
