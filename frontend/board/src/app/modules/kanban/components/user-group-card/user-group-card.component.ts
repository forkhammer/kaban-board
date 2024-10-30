import { Component, Input } from '@angular/core';
import { Group } from '../../models/group';
import { faChevronDown, faChevronUp } from '@fortawesome/free-solid-svg-icons'

@Component({
  selector: 'app-user-group-card',
  templateUrl: './user-group-card.component.html',
  styleUrls: ['./user-group-card.component.scss'],
})
export class UserGroupCardComponent {
  @Input() group!: Group
  @Input() isOpen = true
  @Input() userCount = 0
  faChevronDown = faChevronDown
  faChevronUp = faChevronUp

  toggleOpen() {
    this.isOpen = !this.isOpen
  }
}
