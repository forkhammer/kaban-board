import { HttpErrorResponse } from '@angular/common/http';
import { NotFoundService } from '../services/not-found.service';
import { EMPTY, Observable } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { ToastService } from "../services/toast.service";

export type catchErrorCallback = (err: HttpErrorResponse) => void;

/**
 * Функция pipe обработки ошибок и установки кода ошибки 404 в случае отсутствия объекта
 */
export function catchErrorNotFound<T>(notFound: NotFoundService, callback?: catchErrorCallback) {
  return (source: Observable<T>) =>
    source.pipe(
      catchError(err => {
        if (err.status === 404) {
          notFound.setNotFound();
        }
        if (callback) {
          callback(err);
        }
        return EMPTY;
      }),
    );
}

/**
 * Функция обработки ошибок с выводом сообщения
 */
export function catchErrorMessages<T>(toast: ToastService, callback?: catchErrorCallback) {
  return (source: Observable<T>) =>
    source.pipe(
      catchError(err => {
        if (isValidationError(err)) {
          toast.showMessages(getValidationErrors(err), 'error');
        }

        if (isValidationFormError(err)) {
          toast.showMessages(getValidationFormErrors(err), 'error');
        }

        if (callback) {
          callback(err);
        }
        return EMPTY;
      }),
    );
}

export function isValidationError(err: any): boolean {
  return err.error && Array.isArray(err.error);
}

export function getValidationErrors(err: any): string[] {
  return err.error as string[];
}

export function isValidationFormError(err: any): boolean {
  return err.error && Object.prototype.toString.call(err.error) === '[object Object]';
}

export function getValidationFormErrors(err: any): string[] {
   return Object.keys(err.error).map(key => err.error[key].toString());
}
