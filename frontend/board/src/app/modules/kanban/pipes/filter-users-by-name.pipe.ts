import { Pipe, PipeTransform } from '@angular/core';
import { User } from '../models/user';

@Pipe({
  name: 'filterUsersByName'
})
export class FilterUsersByNamePipe implements PipeTransform {

  transform(users: User[], search: string | null): User[] {
    const lowerText = search ? search.toLowerCase() : null

    return users.filter(user => {
      return lowerText
        ? (user.username.toLowerCase().indexOf(lowerText) > -1) || (user.name.toLowerCase().indexOf(lowerText) > -1)
        : true
    });
  }

}
