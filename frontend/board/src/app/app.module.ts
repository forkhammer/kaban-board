import { NgModule } from '@angular/core';
import { BrowserModule, HammerGestureConfig, HAMMER_GESTURE_CONFIG, HammerModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";
import {HTTP_INTERCEPTORS, HttpClientModule} from "@angular/common/http";
import {CoreModule} from "./modules/core/core.module";
import {UiModule} from "./modules/ui/ui.module";
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import { IndexPageComponent } from './components/index-page/index-page.component';
import {environment} from "../environments/environment";
import {KanbanModule} from "./modules/kanban/kanban.module";
import {JWTInterceptor} from "./modules/core/interceptors/jwt.interceptor";
import { AuthPageComponent } from './components/auth-page/auth-page.component';
import {BootstrapUiModule} from "./modules/bootstrap-ui/bootstrap-ui.module";
import * as Hammer from "hammerjs";


export class MyHammerConfig extends HammerGestureConfig {
  override overrides = <any> {
    swipe: { direction: Hammer.DIRECTION_ALL, velocity: 0.4, threshold: 20, },
  };
}


@NgModule({
  declarations: [
    AppComponent,
    IndexPageComponent,
    AuthPageComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    HttpClientModule,
    CoreModule.forRoot({
      apiUrl: environment.apiUrl
    }),
    UiModule,
    KanbanModule,
    FormsModule,
    ReactiveFormsModule,
    BootstrapUiModule,
    HammerModule,
  ],
  providers: [
    {
      provide: HTTP_INTERCEPTORS,
      useClass: JWTInterceptor,
      multi: true,
    },
    {
      provide: HAMMER_GESTURE_CONFIG,
      useClass: MyHammerConfig,
    },
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
