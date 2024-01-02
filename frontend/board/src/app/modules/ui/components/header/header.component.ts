import {Component, OnDestroy, OnInit} from '@angular/core';
import {GitlabSyncService} from "../../../kanban/services/gitlab-sync.service";
import {distinctUntilChanged, interval, Subject} from "rxjs";
import {map, takeUntil} from "rxjs/operators";
import {ThemeServiceService} from "../../services/theme-service.service";
import {FormBuilder, FormGroup} from "@angular/forms";
import { ActivatedRoute, Router } from '@angular/router';
import { faXmark } from '@fortawesome/free-solid-svg-icons'
import { ClientSettings } from 'src/app/modules/kanban/models/settings';
import { KanbanSettingsService } from 'src/app/modules/kanban/services/kanban-settings.service';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent implements OnDestroy, OnInit {
  public now: Date = new Date()
  private destroy$ = new Subject()
  public form: FormGroup
  public settings: ClientSettings | null = null

  faXmark = faXmark

  constructor(
    public syncService: GitlabSyncService,
    private themeService: ThemeServiceService,
    public fb: FormBuilder,
    private router: Router,
    private route: ActivatedRoute,
    private settingsService: KanbanSettingsService
  ) {
    interval(1000).pipe(
      takeUntil(this.destroy$)
    ).subscribe(_ => {
      this.now = new Date()
    })

    this.form = this.fb.group({
      darkMode: [false],
      search: [''],
    })
  }

  ngOnInit() {
    const search$ = this.route.queryParams.pipe(
      map(params => params['search'])
    );

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

    search$.pipe(
      distinctUntilChanged(),
      takeUntil(this.destroy$)
    ).subscribe(val => {
      this.form.patchValue({search: val})
    })

    this.form.get('search')?.valueChanges.pipe(
      distinctUntilChanged(),
      takeUntil(this.destroy$)
    ).subscribe(val => {
      this.router.navigate([], {queryParams:{search:val}, queryParamsHandling: 'merge'})
    })

    this.settingsService.getClientSettings().pipe(
      takeUntil(this.destroy$)
    ).subscribe(data => {
      this.settings = data
    })
  }

  ngOnDestroy() {
    this.destroy$.next(null)
    this.destroy$.complete()
  }

  clearSearch() {
    this.form.patchValue({search:''})
  }
}
