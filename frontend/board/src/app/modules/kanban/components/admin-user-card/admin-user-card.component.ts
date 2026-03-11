import {Component, DestroyRef, inject, Input, OnInit} from '@angular/core';
import { faEye, faEyeSlash, faChevronDown, faChevronUp } from '@fortawesome/free-solid-svg-icons'
import {ACCOUNT_ROLE_VALUES, User} from "../../models/user";
import {UserService} from "../../services/user.service";
import {distinctUntilChanged, switchMap} from "rxjs/operators";
import { FormBuilder, FormGroup } from '@angular/forms';
import { GroupService } from '../../services/group.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Component({
    selector: 'app-admin-user-card',
    templateUrl: './admin-user-card.component.html',
    styleUrls: ['./admin-user-card.component.scss'],
    standalone: false
})
export class AdminUserCardComponent implements OnInit {
  private userService = inject(UserService)
  private fb = inject(FormBuilder)
  public groupService = inject(GroupService)
  private destoryRef = inject(DestroyRef)

  @Input() user!: User
  form: FormGroup
  isEdit = false

  faEye = faEye
  faEyeSlash = faEyeSlash
  faChevronDown = faChevronDown
  faChevronUp = faChevronUp
  readonly ACCOUNT_ROLE_VALUES = ACCOUNT_ROLE_VALUES

  constructor() {
    this.form = this.fb.group({
      groups: [[]],
      role: [null],
    })
  }

  ngOnInit(): void {
    this.form.patchValue(this.getFormData(this.user), {emitEvent: false})

    this.form.get('groups')?.valueChanges.pipe(
      distinctUntilChanged((x, y) => this.arrayEquals(x, y)),
      switchMap(data => this.userService.setGroups(this.user.id, data)),
      takeUntilDestroyed(this.destoryRef)
    ).subscribe(user => {
      Object.assign(this.user, user)
      this.form.patchValue(this.getFormData(this.user), {emitEvent: false})
    })

    this.form.get('role')?.valueChanges.pipe(
      distinctUntilChanged(),
      switchMap(data => this.userService.setRole(this.user.id, data)),
      takeUntilDestroyed(this.destoryRef)
    ).subscribe(user => {
      Object.assign(this.user, user)
      this.form.patchValue(this.getFormData(this.user), {emitEvent: false})
    })
  }


  toggleVisible() {
    this.userService.setVisibility(this.user.id, !this.user.is_visible).pipe(
      takeUntilDestroyed(this.destoryRef)
    ).subscribe(user => {
      Object.assign(this.user, user)
    })
  }

  toggleEdit() {
    this.isEdit = !this.isEdit
  }

  getFormData(user: User) {
    return {
      groups: user.groups.map(group => group.id),
      role: user.account?.role ?? null,
    }
  }

  private arrayEquals(a: any[], b: any[]) {
    return Array.isArray(a) && Array.isArray(b) && a.length === b.length && a.every((val, index) => val === b[index]);
  }

}
