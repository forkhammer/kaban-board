import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  NgbDropdownModule, NgbModalModule, NgbNavModule, NgbToastModule,
  NgbTooltipModule,
} from '@ng-bootstrap/ng-bootstrap';

@NgModule({
  declarations: [],
  imports: [
    CommonModule,
    NgbToastModule,
    NgbDropdownModule,
    NgbModalModule,
    NgbNavModule,
    NgbTooltipModule
  ],
  exports: [
    NgbToastModule,
    NgbDropdownModule,
    NgbModalModule,
    NgbNavModule,
    NgbTooltipModule,
  ],
})
export class BootstrapUiModule {}
