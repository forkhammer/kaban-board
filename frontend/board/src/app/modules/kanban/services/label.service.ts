import {Injectable, Injector} from '@angular/core';
import {BaseService} from "../../core/services/base.service";
import {Label} from "../models/kanban-label";

@Injectable({
  providedIn: 'root'
})
export class LabelService extends BaseService<Label>{
  public override usePagination = false

  constructor(protected override injector: Injector) {
    super(injector)
    this.apiUrl = this.config.apiUrl + '/labels'
  }
}
