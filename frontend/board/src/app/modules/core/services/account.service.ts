import { Injectable, Injector } from '@angular/core';
import {
  RegistrationRequest,
  RegistrationResult,
  Account,
  AccountAuthResult
} from '../models/account';
import { HttpErrorResponse } from '@angular/common/http';
import { BehaviorSubject, EMPTY, Observable, of } from 'rxjs';
import { Subject } from 'rxjs';
import { takeUntil, switchMap, map, catchError, tap, shareReplay } from 'rxjs/operators';
import { JWTResponse } from '../models/jwt';
import { JWTService } from './jwt.service';
import { BaseService } from './base.service';

@Injectable({
  providedIn: 'root',
})
export class AccountService extends BaseService<Account> {
  user$ = new BehaviorSubject<Account | null>(null);
  isAdmin$ = new BehaviorSubject<boolean>(false);
  userObservable: Observable<Account | null>;
  updateSignal$ = new Subject<any>();
  protected tokenUrl;
  private jwt: JWTService;

  constructor(
    protected override injector: Injector,
  ) {
    super(injector);
    this.jwt = injector.get(JWTService);

    this.apiUrl = this.config.apiUrl + '/account/user';
    this.tokenUrl = this.config.apiUrl + '/account/login';

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
        takeUntil(this.destroy$)
      );

    this.userObservable.subscribe((data: Account | null) => {
      this.user$.next(data);
    });

    this.user$.subscribe(user => {
      this.isAdmin$.next(user != null)
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
}
