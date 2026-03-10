import { DestroyRef, inject, Injectable, Injector } from '@angular/core';
import {
  RegistrationRequest,
  RegistrationResult,
  Account,
  AccountAuthResult,
  AccountRole
} from '../models/account';
import { HttpErrorResponse } from '@angular/common/http';
import { BehaviorSubject, EMPTY, Observable, of } from 'rxjs';
import { Subject } from 'rxjs';
import { takeUntil, switchMap, map, catchError, tap, shareReplay } from 'rxjs/operators';
import { JWTResponse } from '../models/jwt';
import { JWTService } from './jwt.service';
import { BaseService } from './base.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { GitlabAuthCallbackResponse, GitlabAuthResponse } from '../models/gitlab';
import { ToastService } from './toast.service';

@Injectable({
  providedIn: 'root',
})
export class AccountService extends BaseService<Account> {
  user$ = new BehaviorSubject<Account | null>(null);
  isAdmin$ = new BehaviorSubject<boolean>(false);
  userObservable: Observable<Account | null>;
  updateSignal$ = new Subject<any>();
  protected tokenUrl: string
  protected loginGitlabUrl: string
  protected callbackGitlabUrl: string
  private jwt: JWTService

  destroyRef = inject(DestroyRef)
  toast = inject(ToastService)

  constructor(
    protected override injector: Injector,
  ) {
    super(injector);
    this.jwt = injector.get(JWTService);

    this.apiUrl = this.config.apiUrl + '/account/user';
    this.tokenUrl = this.config.apiUrl + '/account/login';
    this.loginGitlabUrl = this.config.apiUrl + '/auth/gitlab';
    this.callbackGitlabUrl = this.config.apiUrl + '/auth/gitlab/callback';

    this.userObservable = this.updateSignal$
      .pipe(
        switchMap(data => {
          return this.getActive().pipe(
            catchError((err: HttpErrorResponse) => {
              return of({user: null});
            }),
          );
        }),
        map(data => data.user),
        takeUntilDestroyed()
      );

    this.userObservable.subscribe((data: Account | null) => {
      this.user$.next(data);
    });

    this.user$.subscribe(user => {
      this.isAdmin$.next(user?.role == AccountRole.ADMIN)
    })

    this.update();
  }

  update() {
    this.updateSignal$.next(true);
  }

  login(username: string, password: string): Observable<AccountAuthResult> {
    return this.http.post(this.tokenUrl, {username, password})
      .pipe(
        switchMap(data => {
          this.jwt.setTokens(data as JWTResponse);
          return this.getActive();
        }),
        tap(data => {
          if (data.user) {
            this.user$.next(data.user);
          }
        }),
        catchError((err: HttpErrorResponse) => {
          switch (err.status) {
            case 400:
              const response = {
                user: null,
                result: false,
                message: err.error.error,
                errors: [],
              };
              // if ('username' in err.error) {
              //   response.errors.push(err.error.username[0]);
              // }
              // if ('password' in err.error) {
              //   response.errors.push(err.error.password[0]);
              // }
              return [response];
            case 401:
              return [
                {
                  user: null,
                  result: false,
                  message: 'Неправильные логин или пароль',
                  errors: ['Неправильные логин или пароль'],
                },
              ];
            case 402:
              console.log(err);
              return [
                {
                  user: null,
                  result: false,
                  message: err.error.detail as string,
                  errors: [err.error.detail as string],
                  need_activate: true,
                },
              ];
          }
          return EMPTY;
        }),
      );
  }

  loginGitlab() {
    return this.http.get<GitlabAuthResponse>(this.loginGitlabUrl)
      .pipe(
        tap(data => {
          if (!data.enabled) {
            this.toast.show(this.toast.createErrorToast('Авторизация через Gitlab отключена'))
          } else if (data.url) {
            window.location.href = data.url
          }
        }),
        takeUntilDestroyed(this.destroyRef)
      )
  }

  register(data: RegistrationRequest) {
    return this.http.post(`${this.apiUrl}/register/`, data).pipe(
      map(data => data as RegistrationResult)
    );
  }

  logout() {
    this.user$.next(null);
    this.jwt.setTokens({
      token: '',
      // refresh: '',
    });
    this.update();
    this.router.navigate(['/auth']);
  }

  getActive() {
    return this.http.get(this.apiUrl).pipe(
      map(data => data as AccountAuthResult),
      shareReplay(),
    );
  }

  exchangeGitlabCode(code: string) {
    return this.http.get<GitlabAuthCallbackResponse>(`${this.callbackGitlabUrl}?code=${code}`)
      .pipe(
        tap(data => {
          this.jwt.setTokens({ token: data.token });
          this.update();
        })
      )
  }
}
