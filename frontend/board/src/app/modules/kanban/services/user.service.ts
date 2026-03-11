import {Injectable, Injector} from '@angular/core';
import {BaseService} from "../../core/services/base.service";
import {AccountRole, User} from "../models/user";
import {map} from "rxjs/operators";

@Injectable({
  providedIn: 'root'
})
export class UserService extends BaseService<User>{
  public override usePagination = false

  constructor(protected override injector: Injector) {
    super(injector)
    console.log('create user service')
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

  setRole(userId: number, role: AccountRole) {
    return this.http.post(`${this.apiUrl}/${userId}/role`, {role}).pipe(
      map(data => data as User)
    )
  }

  formatter(item: User) {
    return item.name
  }
}
