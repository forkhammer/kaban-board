import { Component, forwardRef} from '@angular/core';
import { NG_VALUE_ACCESSOR } from '@angular/forms';
import { SelectComponent } from "../select/select.component";


@Component({
    selector: 'app-select-toggle',
    templateUrl: './select-toggle.component.html',
    styleUrls: ['./select-toggle.component.scss'],
    providers: [
        { provide: NG_VALUE_ACCESSOR, useExisting: forwardRef(() => SelectToggleComponent), multi: true },
    ],
    standalone: false
})
export class SelectToggleComponent extends SelectComponent {

}
