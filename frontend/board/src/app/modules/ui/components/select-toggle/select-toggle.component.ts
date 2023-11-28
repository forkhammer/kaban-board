import { Component, OnInit, forwardRef, Input, OnDestroy, ViewChild, ElementRef } from '@angular/core';
import { ControlValueAccessor, NG_VALUE_ACCESSOR, FormBuilder, FormGroup } from '@angular/forms';
import { faAngleDown, faAngleUp, faTimes } from '@fortawesome/free-solid-svg-icons';
import { Subject } from 'rxjs';
import { SelectValue } from '../../models/select-value';
import { SelectComponent } from "../select/select.component";


@Component({
  selector: 'app-select-toggle',
  templateUrl: './select-toggle.component.html',
  styleUrls: ['./select-toggle.component.scss'],
  providers: [
    { provide: NG_VALUE_ACCESSOR, useExisting: forwardRef(() => SelectToggleComponent), multi: true },
  ]
})
export class SelectToggleComponent extends SelectComponent {

}
