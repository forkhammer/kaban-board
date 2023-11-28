import { Component } from '@angular/core';
import { ToastService } from "../../services/toast.service";

@Component({
  selector: 'app-toast-container',
  templateUrl: './toast-container.component.html',
  styleUrls: ['./toast-container.component.scss']
})
export class ToastContainerComponent {

  constructor(
    public toastService: ToastService
  ) { }

}
