import {Component, inject} from '@angular/core';
import {GitlabSyncService} from "../../../kanban/services/gitlab-sync.service";
import {distinctUntilChanged, interval} from "rxjs";
import {FormBuilder, FormGroup} from "@angular/forms";
import { faXmark } from '@fortawesome/free-solid-svg-icons'
import { ClientSettings } from 'src/app/modules/kanban/models/settings';
import { KanbanSettingsService } from 'src/app/modules/kanban/services/kanban-settings.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { ThemeServiceService } from 'src/app/modules/ui/services/theme-service.service';

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
    })

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

    this.settingsService.getClientSettings().pipe(
      takeUntilDestroyed()
    ).subscribe(data => {
      this.settings = data
    })
  }
}
