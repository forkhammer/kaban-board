import {ModuleWithProviders, NgModule, Optional, SkipSelf} from '@angular/core';
import { CommonModule } from '@angular/common';
import {ToastContainerComponent} from "./components/toast-container/toast-container.component";
import {BootstrapUiModule} from "../bootstrap-ui/bootstrap-ui.module";
import {CoreConfig, CoreConfigService} from "./config";



@NgModule({
  declarations: [
    ToastContainerComponent
  ],
  exports: [
    ToastContainerComponent
  ],
  imports: [
    CommonModule,
    BootstrapUiModule
  ]
})
export class CoreModule {
  constructor(@Optional() @SkipSelf() parentModule: CoreModule) {
    if (parentModule) {
      throw new Error('CoreModule is already loaded. Import it in the AppModule only');
    }
  }

  static forRoot(config: CoreConfig): ModuleWithProviders<any> {
    return {
      ngModule: CoreModule,
      providers: [
        {
          provide: CoreConfigService,
          useValue: config,
        },
      ],
    };
  }
}
