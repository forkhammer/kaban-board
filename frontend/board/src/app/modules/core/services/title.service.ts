import { Injectable, OnDestroy } from '@angular/core';
import { Title, Meta } from '@angular/platform-browser';
import { Router } from '@angular/router';
import { takeUntil } from 'rxjs/operators';
import { BehaviorSubject, Subject, combineLatest } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class TitleService implements OnDestroy {
  public titleSuffix$ = new BehaviorSubject<string>('');
  private titleText$ = new BehaviorSubject<string>('');
  private descriptionText$ = new BehaviorSubject<string>('');
  private destroy$ = new Subject();

  get suffix() {
    return this.titleSuffix$.value;
  }

  set suffix(value: string) {
    this.titleSuffix$.next(value);
  }

  constructor(private title: Title, private meta: Meta, private router: Router) {
    combineLatest([this.titleText$, this.titleSuffix$])
      .pipe(takeUntil(this.destroy$))
      .subscribe(data => {
        let t = data[0];
        if (data[1]) {
          t += ' | ' + data[1];
        }
        this.title.setTitle(t);
        this.meta.updateTag({ property: 'og:title', content: t });
      });

    this.descriptionText$.subscribe(data => {
      this.meta.updateTag({ name: 'description', content: data });
      this.meta.updateTag({ property: 'og:description', content: data });
    });
  }

  ngOnDestroy() {
    this.destroy$.next(null);
    this.destroy$.complete();
  }

  setTitle(title: string) {
    this.titleText$.next(title);
  }

  setDescription(description: string) {
    this.descriptionText$.next(description);
  }

  setTitleAndDescription(title: string) {
    this.setTitle(title)
    this.setDescription(title);
  }
}
