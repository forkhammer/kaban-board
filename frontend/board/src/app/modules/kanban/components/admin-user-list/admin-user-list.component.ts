import {Component, inject} from '@angular/core';
import {finalize} from "rxjs";
import {User} from "../../models/user";
import {UserService} from "../../services/user.service";
import { FormBuilder, FormGroup } from '@angular/forms';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Component({
    selector: 'app-admin-user-list',
    templateUrl: './admin-user-list.component.html',
    styleUrls: ['./admin-user-list.component.scss'],
    standalone: false
})
export class AdminUserListComponent {
  private userService = inject(UserService)
  private fb = inject(FormBuilder)

  public users: User[] = []
  public isLoading = true
  public filterForm: FormGroup

  constructor() {
    this.filterForm = this.fb.group({
      search: [''],
    })
    this.userService.all().pipe(
      finalize(() => this.isLoading = false),
      takeUntilDestroyed()
    ).subscribe(data => {
      this.users = data as User[]
    })
  }

  trackByUser(_: number, user: User): number {
    return user.id
  }

}
