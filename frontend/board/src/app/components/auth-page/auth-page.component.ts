import {Component, DestroyRef, inject} from '@angular/core';
import {FormBuilder, FormGroup, Validators} from "@angular/forms";
import {TitleService} from "../../modules/core/services/title.service";
import {AccountService} from "../../modules/core/services/account.service";
import {Router} from "@angular/router";
import {takeUntilDestroyed} from '@angular/core/rxjs-interop';

@Component({
    selector: 'app-auth-page',
    templateUrl: './auth-page.component.html',
    styleUrls: ['./auth-page.component.scss'],
    standalone: false
})
export class AuthPageComponent {
  title = inject(TitleService);
  accountService = inject(AccountService)
  fb = inject(FormBuilder)
  router = inject(Router)
  destroyRef = inject(DestroyRef)

  form: FormGroup;
  isLoading = false;
  authErrorMessage: string = '';

  constructor() {
    this.title.setTitle('Auth');
    this.form = this.fb.group({
      username: ['', [Validators.required]],
      password: ['', Validators.required],
    });
  }

  submit(e: SubmitEvent) {
    this.isLoading = true;
    this.authErrorMessage = '';

    this.accountService
      .login(this.form.get('username')?.value, this.form.get('password')?.value)
      .pipe(takeUntilDestroyed(this.destroyRef))
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
