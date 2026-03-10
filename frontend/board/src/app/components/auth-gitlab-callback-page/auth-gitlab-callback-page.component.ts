import { Component, inject } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Router } from '@angular/router';
import { AccountService } from '../../modules/core/services/account.service';
import { catchError, EMPTY, filter, map, switchMap, tap } from 'rxjs';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { ToastService } from 'src/app/modules/core/services/toast.service';

@Component({
  selector: 'app-auth-gitlab-callback-page',
  templateUrl: './auth-gitlab-callback-page.component.html',
  styleUrl: './auth-gitlab-callback-page.component.scss',
  standalone: false,
})
export class AuthGitlabCallbackPageComponent {
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private accountService = inject(AccountService);
  private toast = inject(ToastService)

  constructor() {
    this.route.queryParams.pipe(
      map(params => params['code']),
      filter(data => Boolean(data)),
      switchMap(code => {
        return this.accountService.exchangeGitlabCode(code)
      }),
      catchError(err => {
        this.toast.show(this.toast.createErrorToast('Ошибка авторизации'))
        this.router.navigate(['/auth'])
        return EMPTY
      }),
      takeUntilDestroyed()
    ).subscribe(_ => {
      this.router.navigate(['/'])
    })
  }
}
