import {Component, inject, Input} from '@angular/core';
import { ErrorService } from 'src/app/modules/core/services/error.service';

@Component({
  selector: 'app-input-error',
  standalone: false,
  templateUrl: './input-error.component.html',
  styleUrl: './input-error.component.scss'
})
export class InputErrorComponent {
  @Input() errorData: any;
  @Input() code!: string | number;
  errors = inject(ErrorService)
}
