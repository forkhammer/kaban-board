import { InjectionToken } from "@angular/core";

export interface CoreConfig {
  apiUrl: string;
}

export const CoreConfigService = new InjectionToken<CoreConfig>('CoreConfig');
