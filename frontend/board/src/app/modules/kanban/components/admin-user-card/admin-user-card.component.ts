import {Component, Input, OnDestroy, OnInit} from '@angular/core';
import {Subject} from "rxjs";
import { faEye, faEyeSlash, faChevronDown, faChevronUp } from '@fortawesome/free-solid-svg-icons'
import {User} from "../../models/user";
import {UserService} from "../../services/user.service";
import {distinctUntilChanged, switchMap, takeUntil} from "rxjs/operators";
import { FormBuilder, FormGroup } from '@angular/forms';
import { GroupService } from '../../services/group.service';

@Component({
  selector: 'app-admin-user-card',
  templateUrl: './admin-user-card.component.html',
  styleUrls: ['./admin-user-card.component.scss']
})
export class AdminUserCardComponent implements OnDestroy, OnInit {
  @Input() user!: User
  private destroy$ = new Subject()
  form: FormGroup
  isEdit = false

  faEye = faEye
  faEyeSlash = faEyeSlash
  faChevronDown = faChevronDown
  faChevronUp = faChevronUp

  constructor(
    private userService: UserService,
    private fb: FormBuilder,
    public groupService: GroupService
  ) {
    this.form = this.fb.group({
      groups: [[]],
    })
  }

  ngOnInit(): void {
      this.form.patchValue(this.getFormData(this.user))

      this.form.get('groups')?.valueChanges.pipe(
        distinctUntilChanged((x, y) => this.arrayEquals(x, y)),
        switchMap(data => this.userService.setGroups(this.user.id, data)),
        takeUntil(this.destroy$)
      ).subscribe(user => {
        Object.assign(this.user, user)
        this.form.patchValue(this.getFormData(this.user))
      })
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

  toggleEdit() {
    this.isEdit = !this.isEdit
  }

  getFormData(user: User) {
    return {
      groups: user.groups.map(group => group.id)
    }
  }

  private arrayEquals(a: any[], b: any[]) {
    return Array.isArray(a) && Array.isArray(b) && a.length === b.length && a.every((val, index) => val === b[index]);
  }

}
