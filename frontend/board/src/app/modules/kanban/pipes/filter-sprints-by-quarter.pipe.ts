import { Pipe, PipeTransform } from '@angular/core';
import { Sprint } from '../models/sprint';
import { Quarter } from '../models/quarter';

@Pipe({
  name: 'filterSprintsByQuarter',
  standalone: false
})

export class FilterSprintsByQuarterPipe implements PipeTransform {
  transform(sprints: Sprint[], quarter: Quarter): Sprint[] {
    return sprints.filter(s => s.quarter.id === quarter.id)
  }
}
