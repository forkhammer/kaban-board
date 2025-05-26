import { Component, inject } from '@angular/core';
import {ThemeServiceService} from "./modules/ui/services/theme-service.service";

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.scss'],
    standalone: false
})
export class AppComponent {
  public themeService = inject(ThemeServiceService)
}
