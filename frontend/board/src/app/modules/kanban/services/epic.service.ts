import {Injectable, Injector} from '@angular/core';
import {BaseService} from "../../core/services/base.service";
import { Epic } from '../models/epic';

@Injectable({
  providedIn: 'root'
})
export class EpicService extends BaseService<Epic>{
  public override usePagination = false

  constructor(protected override injector: Injector) {
    super(injector)
    this.apiUrl = this.config.apiUrl + '/epic'
  }
}
