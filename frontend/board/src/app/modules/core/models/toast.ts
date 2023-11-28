import { TemplateRef } from "@angular/core";

export interface Toast {
  title?: string;
  text?: string;
  template?: TemplateRef<TemplateRef<any>>;
  className?: string;
  delay?: number;
}
