import { Component, HostBinding, Input } from '@angular/core';

@Component({
    selector: 'app-spinner',
    templateUrl: './spinner.component.html',
    styleUrls: ['./spinner.component.scss'],
    standalone: false
})
export class SpinnerComponent {

  @Input()
  @HostBinding('class.visible')
  visible = false;

  constructor() { }
}
