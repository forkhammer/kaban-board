import {Component, OnDestroy, OnInit} from '@angular/core';
import {GitlabSyncService} from "../../../kanban/services/gitlab-sync.service";
import {distinctUntilChanged, interval, Subject} from "rxjs";
import {takeUntil} from "rxjs/operators";
import {ThemeServiceService} from "../../services/theme-service.service";
import {FormBuilder, FormGroup} from "@angular/forms";

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent implements OnDestroy, OnInit {
  public now: Date = new Date()
  private destroy$ = new Subject()
  public form: FormGroup

  constructor(
    public syncService: GitlabSyncService,
    private themeService: ThemeServiceService,
    public fb: FormBuilder
  ) {
    interval(1000).pipe(
      takeUntil(this.destroy$)
    ).subscribe(_ => {
      this.now = new Date()
    })

    this.form = this.fb.group({
      darkMode: [false]
    })
  }

  ngOnInit() {
    this.themeService.theme$.pipe(
      distinctUntilChanged(),
      takeUntil(this.destroy$)
    ).subscribe(val => {
      this.form.patchValue({darkMode: val === 'dark'})
    })

    this.form.get('darkMode')?.valueChanges.pipe(
      distinctUntilChanged(),
      takeUntil(this.destroy$)
    ).subscribe(value => {
      this.themeService.theme$.next(value ? 'dark' : 'light')
    })
  }

  ngOnDestroy() {
    this.destroy$.next(null)
    this.destroy$.complete()
  }
}
