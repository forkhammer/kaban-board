import { Component, ElementRef, forwardRef, inject, ViewChild } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { ControlValueAccessor, FormBuilder, FormGroup, NG_VALUE_ACCESSOR } from '@angular/forms';
import { BehaviorSubject, combineLatest, of } from 'rxjs';
import { catchError, debounceTime, distinctUntilChanged, filter, switchMap } from 'rxjs/operators';
import { Epic } from '../../models/epic';
import { EpicService } from '../../services/epic.service';
import { faTimes } from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'app-select-epic',
  standalone: false,
  templateUrl: './select-epic.component.html',
  styleUrls: [
    '../../../ui/components/select/select.component.scss',
    '../../../ui/components/select/filter-form.scss',
    './select-epic.component.scss',
  ],
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => SelectEpicComponent),
      multi: true
    }
  ]
})
export class SelectEpicComponent implements ControlValueAccessor {
  protected fb = inject(FormBuilder)
  private epicService = inject(EpicService);
  private onChange: any;

  faTimes = faTimes

  @ViewChild('dropdown') dropdown!: ElementRef;

  value$ = new BehaviorSubject<number | null>(null);
  search$ = new BehaviorSubject<string>('');
  opened$ = new BehaviorSubject<boolean>(false);
  epics: Epic[] = [];
  selectedEpic: Epic | null = null;
  searchForm: FormGroup
  @ViewChild('searchInput') searchInput!: ElementRef<HTMLInputElement>;

  get selectValue(): number | null {
    return this.value$.value;
  }

  set selectValue(item: number | null) {
    if (item) {
      this.writeValue(item);
    } else {
      this.writeValue(null);
    }
  }

  constructor() {
    this.searchForm = this.fb.group({
      search: ['']
    });

    this.searchForm.get('search')?.valueChanges.pipe(
      takeUntilDestroyed()
    ).subscribe(data => {
      this.search$.next(data)
    })

    combineLatest([
      this.opened$.pipe(filter(Boolean), distinctUntilChanged()),
      this.search$.pipe(debounceTime(300), distinctUntilChanged()),
    ]).pipe(
      switchMap(([, search]) =>
        this.epicService.list({ search: search || null }).pipe(
          catchError(() => of([]))
        )
      ),
      takeUntilDestroyed()
    ).subscribe((epics: any) => {
      this.epics = Array.isArray(epics) ? epics : (epics?.results ?? []);
    });

    this.value$.pipe(
      distinctUntilChanged(),
      switchMap(id => {
        if (!id) return of(null);
        if (this.selectedEpic?.id === id) return of(this.selectedEpic);
        return this.epicService.get(id).pipe(catchError(() => of(null)));
      }),
      takeUntilDestroyed()
    ).subscribe(epic => {
      this.selectedEpic = epic as Epic | null;
    });
  }

  writeValue(value: number | null) {
    this.value$.next(value);
    if (this.onChange) this.onChange(value);
  }

  registerOnChange(fn: any) {
    this.onChange = fn;
  }

  registerOnTouched(fn: any) {}

  select(epic: Epic | null) {
    const id = epic?.id ?? null;
    this.selectValue = id;
    (this.dropdown as any).close();
  }

  clear(e: MouseEvent) {
    e.preventDefault()
    e.stopPropagation()
    this.selectValue = null
    return false
  }

  onOpenChange(open: boolean) {
    this.opened$.next(open);
    if (open && this.searchInput) {
      setTimeout(() => this.searchInput.nativeElement.focus(), 0);
    }
  }

  formatEpic(epic: Epic): string {
    return `#${epic.iid} ${epic.title} (${epic.project})`;
  }
}
