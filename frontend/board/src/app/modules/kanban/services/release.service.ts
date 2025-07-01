import {Injectable, Injector} from '@angular/core';
import {BaseService} from "../../core/services/base.service";
import { Release } from '../models/release';

@Injectable({
  providedIn: 'root'
})
export class ReleaseService extends BaseService<Release>{
  public override usePagination = false

  constructor(protected override injector: Injector) {
    super(injector)
    this.apiUrl = this.config.apiUrl + '/release'
  }
}
