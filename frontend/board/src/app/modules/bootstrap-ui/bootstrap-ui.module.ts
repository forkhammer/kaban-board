import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  NgbDropdownModule, NgbModalModule, NgbNavModule, NgbToastModule,
} from '@ng-bootstrap/ng-bootstrap';

@NgModule({
  declarations: [],
  imports: [
    CommonModule,
    NgbToastModule,
    NgbDropdownModule,
    NgbModalModule,
    NgbNavModule,
  ],
  exports: [
    NgbToastModule,
    NgbDropdownModule,
    NgbModalModule,
    NgbNavModule,
  ],
})
export class BootstrapUiModule {}
