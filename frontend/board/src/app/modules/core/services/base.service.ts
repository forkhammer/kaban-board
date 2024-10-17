import { Injectable, Injector, PLATFORM_ID, OnDestroy } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { BaseModel, Pagination } from '../models/base';
import { map } from 'rxjs/operators';
import { Router } from '@angular/router';
import { isPlatformServer } from '@angular/common';
import { CoreConfig, CoreConfigService } from '../config';
import { Subject } from 'rxjs';
import { RestQuery } from '../models/rest';

@Injectable()
export class BaseService<T extends BaseModel> implements OnDestroy {
  protected http: HttpClient;
  protected router: Router;
  protected platformId: any;
  protected RESPONSE: any;
  protected config: CoreConfig;
  protected apiUrl: string;
  protected destroy$ = new Subject();
  public usePagination = true;

  constructor(protected injector: Injector) {
    this.http = injector.get(HttpClient);
    this.router = injector.get(Router);
    this.platformId = injector.get(PLATFORM_ID);
    this.RESPONSE = injector.get('RESPONSE', null);
    this.config = injector.get(CoreConfigService);
    this.apiUrl = this.config.apiUrl + '/api/base/';
  }

  ngOnDestroy() {
    this.destroy$.next(null);
    this.destroy$.complete();
  }

  list(query?: RestQuery) {
    return this.http.get(this.apiUrl, { params: this.filterQuery(query) }).pipe(
      map(res => this.usePagination ? res as Pagination<T> : res as T[])
    );
  }

  all(query?: RestQuery) {
    let params = this.filterQuery(query)
    params = params.set('all', 'true');
    return this.http.get(this.apiUrl, { params }).pipe(
      map(res => this.usePagination ? res as Pagination<T> : res as T[])
    );
  }

  get(id: number | string, query?: RestQuery) {
    return this.http.get(`${this.apiUrl}/${id}`, {params: this.filterQuery(query)}).pipe(map(res => res as T));
  }

  setNotFound() {
    if (isPlatformServer(this.platformId)) {
      this.RESPONSE.status(404);
    }
  }

  save(data: any) {
    if (data.id) {
      return this.http.put(`${this.apiUrl}/${data.id}`, data).pipe(map(res => res as T));
    } else {
      return this.http.post(`${this.apiUrl}`, data).pipe(map(res => res as T));
    }
  }

  patch(data: any) {
    if (data.id) {
      return this.http.patch(`${this.apiUrl}/${data.id}`, data).pipe(map(res => res as T));
    } else {
      return null;
    }
  }

  delete(data: T) {
    return this.http.delete(`${this.apiUrl}/${data.id}`).pipe(map(res => res as T));
  }

  /**
   * Фильтрует параметры запроса от всякого мусора
   * @param query
   * @private
   */
  protected filterQuery(query?: RestQuery): HttpParams {
    let params = new HttpParams();

    if (query !== undefined) {
      Object.keys(query).forEach(key => {
        const value = query[key];

        if (value !== null && value !== undefined) {
          params = params.set(key, value);
        } else {
          params = params.set(key, '');
        }
      });
    }

    return params;
  }
}
