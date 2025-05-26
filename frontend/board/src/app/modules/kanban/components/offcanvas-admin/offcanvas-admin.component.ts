import { Component, inject } from '@angular/core';
import {SettingsService} from "../../../ui/services/settings.service";

@Component({
    selector: 'app-offcanvas-admin',
    templateUrl: './offcanvas-admin.component.html',
    styleUrls: ['./offcanvas-admin.component.scss'],
    standalone: false
})
export class OffcanvasAdminComponent {
  public settingsService = inject(SettingsService)
}
