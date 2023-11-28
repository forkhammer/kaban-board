import { Injectable } from '@angular/core';
import { HttpInterceptor, HttpRequest, HttpHandler, HttpEvent, HttpResponse } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { JWTService } from '../services/jwt.service';
import { map, switchMap } from 'rxjs/operators';

@Injectable({
  providedIn: 'root',
})
export class JWTInterceptor implements HttpInterceptor {
  blacklistedRoutes: Array<string | RegExp> = [/token\/refresh/];

  constructor(private jwt: JWTService) {}

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    const handle$ = of(this.jwt.access).pipe(
      switchMap(data => {
        if (data) {
          // if (this.jwt.isTokenExpired() && !this.isBlacklistedRoute(req)) {
          //   return this.jwt.refreshToken().pipe(
          //     switchMap(token => {
          //       return [this.jwt.access];
          //     }),
          //   );
          // } else {
          //   return [data];
          // }
          return [data];
        } else {
          return [null];
        }
      }),
      switchMap(accessToken => {
        let apiReq = req.clone({
          headers: req.headers.set('ngrok-skip-browser-warning', '123456789'),
        });

        if (accessToken && !this.isBlacklistedRoute(req)) {
          apiReq = apiReq.clone({
            headers: req.headers.set('Authorization', this.jwt.authToken).set('ngrok-skip-browser-warning', '123456789'),
          });
        } else {
          apiReq = apiReq.clone();
        }
        return next.handle(apiReq);
      }),
    );
    return handle$.pipe(
      map((e: HttpEvent<any>) => {
        if (e instanceof HttpResponse) {
          switch ((e as HttpResponse<any>).status) {
            case 401:
              console.log('Не авторизован');
              break;
          }
        }
        return e;
      }),
    );
  }

  isBlacklistedRoute(request: HttpRequest<any>): boolean {
    const url = request.url;

    return (
      this.blacklistedRoutes.findIndex(route =>
        typeof route === 'string' ? route === url : route instanceof RegExp ? route.test(url) : false,
      ) > -1
    );
  }
}
