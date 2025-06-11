import { Component, ElementRef, forwardRef, ViewChild } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';
import { BehaviorSubject, distinctUntilChanged } from 'rxjs';

const ESTIMATE_VALUES = [
  1, 2, 4, 6, 8, 12, 16, 20, 24, 28, 32, 40
]

@Component({
  selector: 'app-select-estimate',
  standalone: false,
  templateUrl: './select-estimate.component.html',
  styleUrl: './select-estimate.component.scss',
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => SelectEstimateComponent),
      multi: true
    }
  ]
})
export class SelectEstimateComponent implements ControlValueAccessor {
  private onChange: any;
  value$ = new BehaviorSubject<number | null>(null);
  @ViewChild('dropdown') dropdown!: ElementRef;

  readonly ESTIMATE_VALUES = ESTIMATE_VALUES

  get selectValue(): number | null {
      return this.value$.value;
  }

  set selectValue(item: number | null) {
    this.writeValue(item);
  }

  constructor() {
    this.value$.pipe(
      distinctUntilChanged(),
      takeUntilDestroyed()
    ).subscribe(data => {
      if (this.onChange) {
        this.onChange(data);
      }
    })
  }


  writeValue(value: number | null) {
    this.value$.next(value);
  }

  registerOnChange(fn: any) {
    this.onChange = fn;
  }

  registerOnTouched(fn: any) {

  }

  select(value: number | null) {
    this.value$.next(value);
    (this.dropdown as any).close();
  }
}
