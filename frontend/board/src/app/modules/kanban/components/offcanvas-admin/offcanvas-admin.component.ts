import { Component } from '@angular/core';
import {SettingsService} from "../../../ui/services/settings.service";

@Component({
  selector: 'app-offcanvas-admin',
  templateUrl: './offcanvas-admin.component.html',
  styleUrls: ['./offcanvas-admin.component.scss']
})
export class OffcanvasAdminComponent {
  constructor(
    public settingsService: SettingsService
  ) {
  }
}
