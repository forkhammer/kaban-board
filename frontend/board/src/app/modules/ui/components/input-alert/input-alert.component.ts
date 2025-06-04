import { Component, ElementRef, HostBinding, inject, Input } from '@angular/core';

@Component({
    selector: 'app-input-alert',
    templateUrl: './input-alert.component.html',
    styleUrls: ['./input-alert.component.scss'],
    host: { class: 'fs-paragraph-xs' },
    standalone: false
})
export class InputAlertComponent {
  @Input() type: string = 'danger';

  private element = inject(ElementRef)

  @HostBinding('class')
  get classes() {
    if (this.element.nativeElement.classList) {
      return this.type;
    } else {
      return this.element.nativeElement.classList.value + ' ' + this.type;
    }
  }

}
