import {Injectable, Injector} from '@angular/core';
import {BaseService} from "../../core/services/base.service";
import { Group } from '../models/group';

@Injectable({
  providedIn: 'root'
})
export class GroupService extends BaseService<Group>{
  public override usePagination = false

  constructor(protected override injector: Injector) {
    super(injector)
    this.apiUrl = this.config.apiUrl + '/groups'
  }
}
