import { Component, ElementRef, forwardRef, Input, ViewChild } from '@angular/core';
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';
import { faTimes } from '@fortawesome/free-solid-svg-icons';
import { ISSUE_PRIORITY_VALUES, IssuePriority } from '../../models/kanban-issue';

@Component({
  selector: 'app-priority-select',
  templateUrl: './priority-select.component.html',
  styleUrls: [
    '../../../ui/components/select/select.component.scss',
    '../../../ui/components/select/filter-form.scss',
    './priority-select.component.scss',
  ],
  providers: [
    { provide: NG_VALUE_ACCESSOR, useExisting: forwardRef(() => PrioritySelectComponent), multi: true },
  ],
  standalone: false,
})
export class PrioritySelectComponent implements ControlValueAccessor {
  @Input() useClear: boolean = false;
  @Input() title = '';
  @Input() nullTitle: string | null = null;
  @ViewChild('dropdown') dropdown!: ElementRef;

  readonly values = ISSUE_PRIORITY_VALUES;
  faTimes = faTimes;

  value: IssuePriority | null = null;
  private onChange: any;

  writeValue(value: IssuePriority | null) {
    this.value = value;
    if (this.onChange) {
      this.onChange(this.value);
    }
  }

  registerOnChange(fn: any) {
    this.onChange = fn;
  }

  registerOnTouched(fn: any) {}

  select(e: MouseEvent, priority: IssuePriority) {
    this.writeValue(priority);
    if (this.dropdown) {
      (this.dropdown as any).close();
    }
    return false;
  }

  clear(e: MouseEvent) {
    this.writeValue(null);
    if (this.dropdown) {
      (this.dropdown as any).close();
    }
    return false;
  }
}
