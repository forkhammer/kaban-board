import { Component } from '@angular/core';
import {ThemeServiceService} from "./modules/ui/services/theme-service.service";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  constructor(
    public themeService: ThemeServiceService
  ) {
  }
}
