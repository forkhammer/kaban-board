import { Component, Input } from '@angular/core';
import { Group } from '../../models/group';

@Component({
  selector: 'app-user-group-card',
  templateUrl: './user-group-card.component.html',
  styleUrls: ['./user-group-card.component.scss']
})
export class UserGroupCardComponent {
  @Input() group!: Group;
}
