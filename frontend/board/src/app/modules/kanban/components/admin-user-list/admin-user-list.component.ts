import {Component, OnDestroy, OnInit} from '@angular/core';
import {finalize, Subject} from "rxjs";
import {takeUntil} from "rxjs/operators";
import {User} from "../../models/user";
import {UserService} from "../../services/user.service";
import { FormBuilder, FormGroup } from '@angular/forms';

@Component({
  selector: 'app-admin-user-list',
  templateUrl: './admin-user-list.component.html',
  styleUrls: ['./admin-user-list.component.scss']
})
export class AdminUserListComponent implements OnInit, OnDestroy {
  private destroy$ = new Subject()
  public users: User[] = []
  public isLoading = true
  public filterForm: FormGroup

  constructor(
    private userService: UserService,
    private fb: FormBuilder
  ) {
    this.filterForm = this.fb.group({
      search: [''],
    })
  }

  ngOnInit() {
    this.userService.all().pipe(
      finalize(() => this.isLoading = false),
      takeUntil(this.destroy$)
    ).subscribe(data => {
      this.users = data as User[]
    })
  }

  ngOnDestroy() {
    this.destroy$.next(null)
    this.destroy$.complete()
  }

  trackByUser(_: number, user: User): number {
    return user.id
  }

}
