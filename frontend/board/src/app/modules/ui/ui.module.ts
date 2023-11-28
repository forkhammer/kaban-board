import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HeaderComponent } from './components/header/header.component';
import {SpinnerComponent} from "./components/spinner/spinner.component";
import { HeaderAccountComponent } from './components/header-account/header-account.component';
import {BootstrapUiModule} from "../bootstrap-ui/bootstrap-ui.module";
import {RouterModule} from "@angular/router";
import {InputAlertComponent} from "./components/input-alert/input-alert.component";
import {SelectComponent} from "./components/select/select.component";
import {SelectModelComponent} from "./components/select-model/select-model.component";
import {SelectModelMultipleComponent} from "./components/select-model-multiple/select-model-multiple.component";
import {SelectMultipleComponent} from "./components/select-multiple/select-multiple.component";
import {SelectToggleComponent} from "./components/select-toggle/select-toggle.component";
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import {FontAwesomeModule} from "@fortawesome/angular-fontawesome";
import {SafePipe} from "./pipes/safe.pipe";
import {SelectItemsPipe} from "./pipes/select-items.pipe";
import { ElapsedTimePipe } from './pipes/elapsed-time.pipe';



@NgModule({
  declarations: [
    HeaderComponent,
    SpinnerComponent,
    HeaderAccountComponent,
    InputAlertComponent,
    SelectComponent,
    SelectModelComponent,
    SelectModelMultipleComponent,
    SelectMultipleComponent,
    SelectToggleComponent,
    SafePipe,
    SelectItemsPipe,
    ElapsedTimePipe,
  ],
  exports: [
    HeaderComponent,
    SpinnerComponent,
    InputAlertComponent,
    SelectComponent,
    SelectModelComponent,
    SelectModelMultipleComponent,
    SelectMultipleComponent,
    SelectToggleComponent,
    SafePipe,
    ElapsedTimePipe,
  ],
  imports: [
    CommonModule,
    BootstrapUiModule,
    RouterModule,
    FormsModule,
    ReactiveFormsModule,
    FontAwesomeModule,
  ]
})
export class UiModule { }
