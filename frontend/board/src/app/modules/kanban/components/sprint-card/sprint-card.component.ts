import { Component, Input } from '@angular/core';
import { Sprint, SprintStatus } from '../../models/sprint';
import {faClock} from '@fortawesome/free-regular-svg-icons';
import {faPlay, faStop} from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'app-sprint-card',
  standalone: false,
  templateUrl: './sprint-card.component.html',
  styleUrl: './sprint-card.component.scss'
})
export class SprintCardComponent {
  @Input() sprint!: Sprint

  SprintStatus = SprintStatus
  faClock = faClock
  faPlay = faPlay
  faStop = faStop
}
