import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'formatDate',
  standalone: false
})
export class FormatDatePipe implements PipeTransform {

  transform(value: Date | string | null | undefined, format: 'long' | 'short' = 'long'): Date | string {
    try {
      return value ?
              new Date(value).toLocaleString('ru', {
                year: 'numeric',
                month: format,
                day: 'numeric'
              })
            : '';
    } catch {
      return '';
    }
  }

}
