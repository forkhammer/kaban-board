import { Pipe, PipeTransform } from '@angular/core';
import { SelectValue } from '../models/select-value';

@Pipe({
  name: 'selectItemsFilter',
})
export class SelectItemsPipe implements PipeTransform {
  transform(values: SelectValue[], search: string | null): any {
    if (!search) {
      return values;
    }
    return values.filter(item => {
      return item.title.toLowerCase().indexOf(search.toLowerCase()) > -1;
    });
  }
}
