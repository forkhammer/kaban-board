import { Pipe, PipeTransform } from '@angular/core';
import { Project } from '../models/project';

@Pipe({
  name: 'filterProjectsByName'
})
export class FilterProjectsByNamePipe implements PipeTransform {

  transform(projects: Project[], search: string | null): Project[] {
    const lowerText = search ? search.toLowerCase() : null

    return projects.filter(project => {
      return lowerText
        ? (project.name.toLowerCase().indexOf(lowerText) > -1)
        : true
    });
  }

}
