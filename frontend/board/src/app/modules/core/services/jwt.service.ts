import {Injectable, Inject, PLATFORM_ID} from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import { JWTResponse } from '../models/jwt';
import { HttpClient } from '@angular/common/http';
import {catchError, tap, shareReplay, map} from 'rxjs/operators';
import { CoreConfigService, CoreConfig } from '../config';
import {Observable, of} from 'rxjs';


@Injectable({
  providedIn: 'root',
})
export class JWTService {
  public access = '';
  // public refresh = '';
  protected tokenUrl = this.config.apiUrl + '/account/login';
  private refresh$: Observable<JWTResponse | null> | null  = null;

  constructor(
    private http: HttpClient,
    @Inject(PLATFORM_ID) private platfromId: object,
    @Inject(CoreConfigService) private config: CoreConfig,
  ) {
    this.loadFromStorage();
    this.loadQueryTokens();
  }

  refreshToken(): Observable<JWTResponse | null> {
    console.log('refresh tokens');

    if (this.refresh$) {
      return this.refresh$;
    } else {
      console.log('send request tokens');

      // @ts-ignore
      this.refresh$ = this.http.post(`${this.tokenUrl}refresh/`, { refresh: this.refresh }).pipe(
        map(data => data as JWTResponse),
        catchError(err => {
          console.error('error get tokens');
          return of(null);
        }),
        tap(data => {
          console.log('get tokens');
          this.setTokens(data as JWTResponse);
          this.refresh$ = null;
        }),
        shareReplay(10, 5000),
      );

      return this.refresh$;
    }
  }

  get authToken(): string {
    if (this.access) {
      return 'Bearer ' + this.access;
    } else {
      return '';
    }
  }

  loadFromStorage() {
    if (isPlatformBrowser(this.platfromId)) {
      this.access = localStorage.getItem('accessToken') || '';
      // this.refresh = localStorage.getItem('refreshToken') || '';
    }
  }

  saveToStorage() {
    if (isPlatformBrowser(this.platfromId)) {
      localStorage.setItem('accessToken', this.access);
      // localStorage.setItem('refreshToken', this.refresh);
    }
  }

  setTokens(tokens: JWTResponse | null) {
    if (tokens) {
      this.access = tokens.token;
      // this.refresh = tokens.refresh;
    } else {
      this.access = '';
      // this.refresh = '';
    }
    this.saveToStorage();
  }

  isTokenExpired(offsetSeconds?: number): boolean {
    if (this.access === null || this.access === '') {
      return true;
    }
    let date = this.getTokenExpirationDate();
    offsetSeconds = offsetSeconds || 0;

    if (date === null) {
      return true;
    }

    return !(date.valueOf() > new Date().valueOf() + offsetSeconds * 1000);
  }

  getTokenExpirationDate(): Date | null {
    let decoded: any;
    decoded = this.decodeToken();

    if (!decoded.hasOwnProperty('exp')) {
      return null;
    }

    const date = new Date(0);
    date.setUTCSeconds(decoded.exp);

    return date;
  }

  private decodeToken(): any {
    if (this.access === null) {
      return null;
    }

    const parts = this.access.split('.');

    if (parts.length !== 3) {
      throw new Error(
        "The inspected token doesn't appear to be a JWT. Check to make sure it has three parts and see https://jwt.io for more.",
      );
    }

    const decoded = this.urlBase64Decode(parts[1]);
    if (!decoded) {
      throw new Error('Cannot decode the token.');
    }

    return JSON.parse(decoded);
  }

  private urlBase64Decode(str: string): string {
    let output = str.replace(/-/g, '+').replace(/_/g, '/');
    switch (output.length % 4) {
      case 0: {
        break;
      }
      case 2: {
        output += '==';
        break;
      }
      case 3: {
        output += '=';
        break;
      }
      default: {
        throw new Error('Illegal base64url string!');
      }
    }
    return this.b64DecodeUnicode(output);
  }

  private b64DecodeUnicode(str: any) {
    return decodeURIComponent(
      Array.prototype.map
        .call(this.b64decode(str), (c: any) => {
          return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
        })
        .join(''),
    );
  }

  private b64decode(str: string): string {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=';
    let output = '';

    str = String(str).replace(/=+$/, '');

    if (str.length % 4 === 1) {
      throw new Error("'atob' failed: The string to be decoded is not correctly encoded.");
    }

    for (
      // initialize result and counters
      let bc = 0, bs: any, buffer: any, idx = 0;
      // get next character
      (buffer = str.charAt(idx++));
      // character found in table? initialize bit storage and add its ascii value;
      ~buffer &&
      ((bs = bc % 4 ? bs * 64 + buffer : buffer),
      // and if not first of each 4 characters,
      // convert the first 8 bits to one ascii character
      bc++ % 4)
        ? (output += String.fromCharCode(255 & (bs >> ((-2 * bc) & 6))))
        : 0
    ) {
      // try to find character in table (0-63, not found => -1)
      buffer = chars.indexOf(buffer);
    }
    return output;
  }

  private loadQueryTokens() {
    if (!this.access) {
      const queryParams = this.getQueryParams();
      if (queryParams) {
        const access = queryParams.get('accessToken');
        // const refresh = queryParams.get('refreshToken');
        if (access) {
          this.setTokens({token: access});
        }
      }
    }
  }

  private getQueryParams(): URLSearchParams | null {
    if (isPlatformBrowser(this.platfromId)) {
      return new URLSearchParams(window.location.search);
    }

    return null;
  }
}
