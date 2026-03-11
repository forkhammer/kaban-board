import {Component, DestroyRef, inject} from '@angular/core';
import {FormBuilder, FormGroup, Validators} from "@angular/forms";
import {TitleService} from "../../modules/core/services/title.service";
import {AccountService} from "../../modules/core/services/account.service";
import {Router} from "@angular/router";
import {takeUntilDestroyed} from '@angular/core/rxjs-interop';
import { faSquareGitlab } from '@fortawesome/free-brands-svg-icons';
import { catchError, EMPTY } from 'rxjs';
import {ThemeServiceService} from "../../modules/ui/services/theme-service.service";

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
  themeService = inject(ThemeServiceService)

  faSquareGitlab = faSquareGitlab

  form: FormGroup;
  isLoading = false;
  isLoadingGitlab = false
  authErrorMessage: string = '';
  theme: string = 'light';

  constructor() {
    this.title.setTitle('Auth');
    this.form = this.fb.group({
      username: ['', [Validators.required]],
      password: ['', Validators.required],
    });

    this.themeService.theme$.subscribe(theme => {
      this.theme = theme;
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

  loginGitlab(e: MouseEvent) {
    this.isLoadingGitlab = true
    this.accountService.loginGitlab().pipe(
      catchError(err => {
        this.isLoadingGitlab = false
        return EMPTY
      })
    ).subscribe(data => {
      this.isLoadingGitlab = false
    })
    e.preventDefault()
    return false
  }
}
