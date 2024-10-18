import {Injectable, Injector} from '@angular/core';
import {BaseService} from "../../core/services/base.service";
import {User} from "../models/user";
import {map} from "rxjs/operators";

@Injectable({
  providedIn: 'root'
})
export class UserService extends BaseService<User>{
  public override usePagination = false

  constructor(protected override injector: Injector) {
    super(injector)
    this.apiUrl = this.config.apiUrl + '/users'
  }

  setVisibility(userId: number, visibility: boolean) {
    return this.http.post(`${this.apiUrl}/${userId}/visibility`, {visible: visibility}).pipe(
      map(data => data as User)
    )
  }

  setGroups(userId: number, groups: number[]) {
    return this.http.post(`${this.apiUrl}/${userId}/groups`, {groups}).pipe(
      map(data => data as User)
    )
  }
}
