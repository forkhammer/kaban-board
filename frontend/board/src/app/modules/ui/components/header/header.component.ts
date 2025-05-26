import {Component, inject} from '@angular/core';
import {GitlabSyncService} from "../../../kanban/services/gitlab-sync.service";
import {distinctUntilChanged, interval} from "rxjs";
import {map} from "rxjs/operators";
import {ThemeServiceService} from "../../services/theme-service.service";
import {FormBuilder, FormGroup} from "@angular/forms";
import { ActivatedRoute, Router } from '@angular/router';
import { faXmark } from '@fortawesome/free-solid-svg-icons'
import { ClientSettings } from 'src/app/modules/kanban/models/settings';
import { KanbanSettingsService } from 'src/app/modules/kanban/services/kanban-settings.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Component({
    selector: 'app-header',
    templateUrl: './header.component.html',
    styleUrls: ['./header.component.scss'],
    standalone: false
})
export class HeaderComponent {
  public syncService = inject(GitlabSyncService);
  private themeService =  inject(ThemeServiceService);
  public fb = inject(FormBuilder);
  private router  = inject(Router);
  private route = inject(ActivatedRoute)
  private settingsService = inject(KanbanSettingsService)

  public now: Date = new Date()
  public form: FormGroup
  public settings: ClientSettings | null = null

  faXmark = faXmark

  constructor() {
    interval(1000).pipe(
      takeUntilDestroyed()
    ).subscribe(_ => {
      this.now = new Date()
    })

    this.form = this.fb.group({
      darkMode: [false],
      search: [''],
    })

    const search$ = this.route.queryParams.pipe(
      map(params => params['search'])
    );

    this.themeService.theme$.pipe(
      distinctUntilChanged(),
      takeUntilDestroyed()
    ).subscribe(val => {
      this.form.patchValue({darkMode: val === 'dark'})
    })

    this.form.get('darkMode')?.valueChanges.pipe(
      distinctUntilChanged(),
      takeUntilDestroyed()
    ).subscribe(value => {
      this.themeService.theme$.next(value ? 'dark' : 'light')
    })

    search$.pipe(
      distinctUntilChanged(),
      takeUntilDestroyed()
    ).subscribe(val => {
      this.form.patchValue({search: val})
    })

    this.form.get('search')?.valueChanges.pipe(
      distinctUntilChanged(),
      takeUntilDestroyed()
    ).subscribe(val => {
      this.router.navigate([], {queryParams:{search:val}, queryParamsHandling: 'merge'})
    })

    this.settingsService.getClientSettings().pipe(
      takeUntilDestroyed()
    ).subscribe(data => {
      this.settings = data
    })
  }

  clearSearch() {
    this.form.patchValue({search:''})
  }
}
