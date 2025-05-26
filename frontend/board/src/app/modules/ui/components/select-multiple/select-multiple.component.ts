import { Component, forwardRef, Input, ViewChild, ElementRef, inject } from '@angular/core';
import { ControlValueAccessor, NG_VALUE_ACCESSOR, FormBuilder, FormGroup } from '@angular/forms';
import { faAngleDown, faAngleUp, faTimes } from '@fortawesome/free-solid-svg-icons';
import { SelectValue, SelectValueIdentity } from '../../models/select-value';


@Component({
    selector: 'app-select-multiple',
    templateUrl: './select-multiple.component.html',
    styleUrls: ['./select-multiple.component.scss'],
    providers: [
        { provide: NG_VALUE_ACCESSOR, useExisting: forwardRef(() => SelectMultipleComponent), multi: true },
    ],
    standalone: false
})
export class SelectMultipleComponent implements ControlValueAccessor {
  private fb = inject(FormBuilder)

  @Input() values: SelectValue[] = [];
  @Input() title = '';
  @Input() useSearch: boolean = false;
  @Input() useClear: boolean = false;
  @Input() nullTitle: string | null = null;
  @Input() formatter: any = null;
  @Input() itemFormatter: any = null;
  @ViewChild('dropdown') dropdown!: ElementRef;

  public value: SelectValueIdentity[] = [];
  private onChange: any;
  faAngleDown = faAngleDown;
  faAngleUp = faAngleUp;
  faTimes = faTimes;
  searchForm: FormGroup;

  get selectValue(): SelectValueIdentity[] {
    return this.value;
  }

  set selectValue(item: SelectValueIdentity[]) {
    if (item) {
      this.writeValue(item);
    } else {
      this.writeValue([]);
    }
  }

  constructor(
  ) {
    this.searchForm = this.fb.group({
      search: ['']
    });
  }

  writeValue(value: SelectValueIdentity[]) {
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
    if (this.selectValue.indexOf(item.id) === -1) {
      this.selectValue = this.selectValue.concat(item.id);
    } else {
      this.selectValue = this.selectValue.filter(value => value !== item.id);
    }
    return false;
  }

  clear(e: MouseEvent) {
    this.selectValue = [];
    return false;
  }

  getTitle(item: SelectValue) {
    if (this.formatter) {
      return this.formatter(item);
    }
    return item.title;
  }
}
