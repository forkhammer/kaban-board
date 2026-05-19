import { Component, ElementRef, forwardRef, Input, ViewChild } from '@angular/core';
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';
import { faTimes } from '@fortawesome/free-solid-svg-icons';
import { BIND_STATUS_LABELS, BIND_STATUS_VALUES, BindStatus } from '../../models/kanban-issue';

@Component({
  selector: 'app-binding-status-select',
  templateUrl: './binding-status-select.component.html',
  styleUrls: [
    '../../../ui/components/select/select.component.scss',
    '../../../ui/components/select/filter-form.scss',
    './binding-status-select.component.scss',
  ],
  providers: [
    { provide: NG_VALUE_ACCESSOR, useExisting: forwardRef(() => BindingStatusSelectComponent), multi: true },
  ],
  standalone: false,
})
export class BindingStatusSelectComponent implements ControlValueAccessor {
  @Input() useClear: boolean = false;
  @Input() title = '';
  @Input() nullTitle: string | null = null;
  @ViewChild('dropdown') dropdown!: ElementRef;

  readonly values = BIND_STATUS_VALUES;
  readonly BIND_STATUS_LABELS = BIND_STATUS_LABELS;
  readonly BindStatus = BindStatus;
  faTimes = faTimes;

  value: BindStatus | null = null;
  private onChange: any;

  private emitChange(value: BindStatus | null): void {
    if (this.onChange) {
      this.onChange(value);
    }
  }

  writeValue(value: BindStatus | null) {
    this.value = value;
  }

  registerOnChange(fn: any) {
    this.onChange = fn;
  }

  registerOnTouched(fn: any) {}

  select(e: MouseEvent, status: BindStatus) {
    this.writeValue(status);
    this.emitChange(status);
    if (this.dropdown) {
      (this.dropdown as any).close();
    }
    return false;
  }

  clear(e: MouseEvent) {
    this.writeValue(null);
    this.emitChange(null);
    if (this.dropdown) {
      (this.dropdown as any).close();
    }
    return false;
  }
}
