import {Component, Input, OnDestroy} from '@angular/core';
import {Subject} from "rxjs";
import {ToastService} from "../../../core/services/toast.service";
import { faEye, faEyeSlash } from '@fortawesome/free-solid-svg-icons'
import {User} from "../../models/user";
import {UserService} from "../../services/user.service";
import {takeUntil} from "rxjs/operators";

@Component({
  selector: 'app-admin-user-card',
  templateUrl: './admin-user-card.component.html',
  styleUrls: ['./admin-user-card.component.scss']
})
export class AdminUserCardComponent implements OnDestroy {
  @Input() user!: User
  private destroy$ = new Subject()

  faEye = faEye
  faEyeSlash = faEyeSlash

  constructor(
    private userService: UserService,
    private toast: ToastService
  ) {
  }

  ngOnDestroy() {
    this.destroy$.next(null)
    this.destroy$.complete()
  }

  toggleVisible() {
    this.userService.setVisibility(this.user.id, !this.user.is_visible).pipe(
      takeUntil(this.destroy$)
    ).subscribe(user => {
      Object.assign(this.user, user)
    })
  }
}
