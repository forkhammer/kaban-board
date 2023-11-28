import {Component, OnDestroy, OnInit} from '@angular/core';
import {Subject} from "rxjs";
import {FormBuilder, FormGroup, Validators} from "@angular/forms";
import {TitleService} from "../../modules/core/services/title.service";
import {AccountService} from "../../modules/core/services/account.service";
import {Router} from "@angular/router";
import {takeUntil} from "rxjs/operators";

@Component({
  selector: 'app-auth-page',
  templateUrl: './auth-page.component.html',
  styleUrls: ['./auth-page.component.scss']
})
export class AuthPageComponent implements OnInit, OnDestroy {
  destroy$ = new Subject();
  form: FormGroup;
  isLoading = false;
  authErrorMessage: string = '';

  constructor(
    private title: TitleService,
    private userService: AccountService,
    private fb: FormBuilder,
    private router: Router,
  ) {
    this.form = this.fb.group({
      username: ['', [Validators.required]],
      password: ['', Validators.required],
    });
  }

  ngOnInit(): void {
    this.title.setTitle('Auth');
  }

  ngOnDestroy() {
    this.destroy$.next(null);
    this.destroy$.complete();
  }

  submit(e: SubmitEvent) {
    this.isLoading = true;
    this.authErrorMessage = '';

    this.userService
      .login(this.form.get('username')?.value, this.form.get('password')?.value)
      .pipe(takeUntil(this.destroy$))
      .subscribe(data => {
        this.isLoading = false;
        if (!data.user) {
          if (data.message) {
            this.authErrorMessage = data.message;
          }
        } else {
          let url = window.location.pathname;
          if (url === '/auth') {
            this.router.navigate(['/']);
          } else {
            this.router.navigate([url]);
          }
        }
      });
    return false;
  }
}
