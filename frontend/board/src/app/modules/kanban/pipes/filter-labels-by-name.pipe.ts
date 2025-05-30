import { Pipe, PipeTransform } from '@angular/core';
import { Label } from '../models/kanban-label';

@Pipe({
    name: 'filterLabelsByName',
    standalone: false
})
export class FilterLabelsByNamePipe implements PipeTransform {

  transform(labels: Label[], search: string | null): Label[] {
    const lowerText = search ? search.toLowerCase() : null

    return labels.filter(label => {
      return lowerText
        ? (label.name.toLowerCase().indexOf(lowerText) > -1)
        : true
    });
  }

}
