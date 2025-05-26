import { Injectable } from '@angular/core';
import { Title, Meta } from '@angular/platform-browser';
import { Router } from '@angular/router';
import { BehaviorSubject, combineLatest } from 'rxjs';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Injectable({
  providedIn: 'root',
})
export class TitleService {
  public titleSuffix$ = new BehaviorSubject<string>('');
  private titleText$ = new BehaviorSubject<string>('');
  private descriptionText$ = new BehaviorSubject<string>('');

  get suffix() {
    return this.titleSuffix$.value;
  }

  set suffix(value: string) {
    this.titleSuffix$.next(value);
  }

  constructor(private title: Title, private meta: Meta, private router: Router) {
    combineLatest([this.titleText$, this.titleSuffix$])
      .pipe(takeUntilDestroyed())
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
