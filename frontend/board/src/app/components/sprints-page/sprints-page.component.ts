import { Component, inject } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { FormBuilder, FormGroup } from '@angular/forms';
import { faArrowLeftLong } from '@fortawesome/free-solid-svg-icons'
import { ToastService } from 'src/app/modules/core/services/toast.service';
import { catchErrorMessages } from 'src/app/modules/core/tools/catch-error';
import { Quarter } from 'src/app/modules/kanban/models/quarter';
import { Sprint } from 'src/app/modules/kanban/models/sprint';
import { SprintService } from 'src/app/modules/kanban/services/sprint.service';
import { TeamService } from 'src/app/modules/kanban/services/team.service';
import {faPlus, faMinus} from '@fortawesome/free-solid-svg-icons'

@Component({
  selector: 'app-sprints-page',
  standalone: false,
  templateUrl: './sprints-page.component.html',
  styleUrl: './sprints-page.component.scss'
})
export class SprintsPageComponent {
  teamService = inject(TeamService)
  fb = inject(FormBuilder)
  sprintService = inject(SprintService)
  toast = inject(ToastService)

  faArrowLeftLong = faArrowLeftLong
  faPlus = faPlus
  faMinus = faMinus

  form: FormGroup
  sprints: Sprint[] = []
  quarters: Quarter[] = []

  constructor() {
    this.form = this.fb.group({
      team: [null],
      quarter: [null]
    })

    this.sprintService.list().pipe(
      catchErrorMessages(this.toast),
      takeUntilDestroyed()
    ).subscribe(data => {
      this.sprints = data as Sprint[]

      this.quarters = this.sprints.reduce((quarters, s) => {
        if (!quarters.find(q => q.id === s.quarter.id)) {
          quarters.push(s.quarter)
        }
        return quarters
      }, [] as Quarter[])
    })
  }
}
