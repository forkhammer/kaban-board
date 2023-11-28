import { Component } from '@angular/core';
import {AccountService} from "../../../core/services/account.service";
import { faGear, faCircleXmark, faUser } from '@fortawesome/free-solid-svg-icons';
import { faCircleUser } from '@fortawesome/free-regular-svg-icons';
import {SettingsService} from "../../services/settings.service";

@Component({
  selector: 'app-header-account',
  templateUrl: './header-account.component.html',
  styleUrls: ['./header-account.component.scss']
})
export class HeaderAccountComponent {
  faGear = faGear
  faXmark = faCircleXmark
  faUser = faCircleUser

  constructor(
    public accountService: AccountService,
    public settingsService: SettingsService
  ) {
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
