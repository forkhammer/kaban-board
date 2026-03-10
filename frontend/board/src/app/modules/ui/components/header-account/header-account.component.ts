import { Component, inject } from '@angular/core';
import {AccountService} from "../../../core/services/account.service";
import { faGear, faCircleXmark } from '@fortawesome/free-solid-svg-icons';
import { faCircleUser } from '@fortawesome/free-regular-svg-icons';
import {SettingsService} from "../../services/settings.service";
import { Account } from 'src/app/modules/core/models/account';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Component({
    selector: 'app-header-account',
    templateUrl: './header-account.component.html',
    styleUrls: ['./header-account.component.scss'],
    standalone: false
})
export class HeaderAccountComponent {
  public accountService = inject(AccountService)
  public settingsService = inject(SettingsService)

  faGear = faGear
  faXmark = faCircleXmark
  faUser = faCircleUser

  account: Account | null = null

  constructor() {
    this.accountService.user$.pipe(takeUntilDestroyed()).subscribe(data => this.account = data)
  }

  logout() {
    this.accountService.logout()
  }

  login(e: MouseEvent) {
    e.preventDefault()
    return false
  }

  toggleSettings() {
    this.settingsService.toggle()
  }

  openSettings() {
    this.settingsService.open()
  }
}
