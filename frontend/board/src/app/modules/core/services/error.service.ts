import { Injectable } from '@angular/core';
import {isArray} from "lodash";

@Injectable({
  providedIn: 'root'
})
export class ErrorService {

  getError(error: any, code: string | number): string | undefined {
    if (error && error.hasOwnProperty(code)) {
      const errors = error[code];
      if (isArray(errors)) {
        return errors[0];
      } else {
        return errors;
      }
    }

    return undefined;
  }
}
