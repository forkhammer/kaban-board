import { inject, Injectable, Injector, PLATFORM_ID } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { BaseModel, Pagination } from '../models/base';
import { map, tap } from 'rxjs/operators';
import { Router } from '@angular/router';
import { isPlatformServer } from '@angular/common';
import { CoreConfigService } from '../config';
import { RestQuery } from '../models/rest';
import { CollectionCache } from './collection-cache';
import { of } from 'rxjs';

@Injectable()
export class BaseService<T extends BaseModel> {
  protected http = inject(HttpClient);
  protected router = inject(Router)
  protected platformId = inject(PLATFORM_ID);
  protected RESPONSE: any;
  protected config =  inject(CoreConfigService);
  protected apiUrl: string;
  public usePagination = true;
  protected useCache = true;
  protected cache: CollectionCache<T>;

  constructor(protected injector: Injector) {
    this.RESPONSE = injector.get('RESPONSE', null);
    this.apiUrl = this.config.apiUrl + '/api/base/';
    this.cache = new CollectionCache<T>()
  }

  list(query?: RestQuery) {
    return this.http.get(this.apiUrl, { params: this.filterQuery(query) }).pipe(
      map(res => this.usePagination ? res as Pagination<T> : res as T[]),
      tap(data => {
        if (this.useCache) {
          this.cache.setItems(this.usePagination ? (data as Pagination<T>).results : (data as T[]));
        }
      })
    );
  }

  all(query?: RestQuery) {
    let params = this.filterQuery(query)
    params = params.set('all', 'true');
    return this.http.get(this.apiUrl, { params }).pipe(
      map(res => this.usePagination ? res as Pagination<T> : res as T[]),
      tap(data => {
        if (this.useCache) {
          this.cache.setItems(this.usePagination ? (data as Pagination<T>).results : (data as T[]));
        }
      })
    );
  }

  get(id: number | string, query?: RestQuery) {
    if (this.useCache) {
      const item = this.cache.get(id);
      if (item) {
        return of(item);
      }
    }
    return this.http.get(`${this.apiUrl}/${id}`, {params: this.filterQuery(query)}).pipe(
      map(res => res as T),
      tap(this.updateItemCache.bind(this)),
    );
  }

  setNotFound() {
    if (isPlatformServer(this.platformId)) {
      this.RESPONSE.status(404);
    }
  }

  save(data: any) {
    if (data.id) {
      return this.http.put(`${this.apiUrl}/${data.id}`, data).pipe(
        map(res => res as T),
        tap(this.updateItemCache.bind(this)),
      );
    } else {
      return this.http.post(`${this.apiUrl}`, data).pipe(
        map(res => res as T),
        tap(this.updateItemCache.bind(this)),
      );
    }
  }

  patch(data: any) {
    if (data.id) {
      return this.http.patch(`${this.apiUrl}/${data.id}`, data).pipe(
        map(res => res as T),
        tap(this.updateItemCache.bind(this)),
      );
    } else {
      return null;
    }
  }

  delete(data: T) {
    return this.http.delete(`${this.apiUrl}/${data.id}`).pipe(
      map(res => res as T),
      tap(this.updateItemCache.bind(this)),
    );
  }

  invalidateCache() {
    if (this.useCache) {
      this.cache.invalidate();
    }
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

  protected updateItemCache(data: T) {
    if (this.useCache) {
      this.cache.setItems([data]);
    }
  }
}
