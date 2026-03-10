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
import { ActivatedRoute, Router } from '@angular/router';
import { BehaviorSubject, combineLatestWith, debounceTime, distinctUntilChanged, map, switchMap } from 'rxjs';
import { SprintModalServiceService } from 'src/app/modules/kanban/services/sprint-modal.service';
import { SelectValue } from 'src/app/modules/ui/models/select-value';
import { TitleService } from 'src/app/modules/core/services/title.service';
import { AccountService } from 'src/app/modules/core/services/account.service';

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
  router = inject(Router)
  route = inject(ActivatedRoute)
  sprintModal = inject(SprintModalServiceService)
  title = inject(TitleService)
  accountService = inject(AccountService)

  faArrowLeftLong = faArrowLeftLong
  faPlus = faPlus
  faMinus = faMinus

  form: FormGroup
  sprints: Sprint[] = []
  quarters: Quarter[] = []
  allQuarters: Quarter[] = []
  quartersValues: SelectValue[] = []
  reload$ = new BehaviorSubject<null>(null)

  constructor() {
    this.title.setTitleAndDescription('Sprints')

    this.form = this.fb.group({
      team: [null],
      quarter: [null]
    })

    const team$ = this.route.queryParams.pipe(
      map(params => params['team'] ? Number(params['team']) : null),
      distinctUntilChanged()
    );
    const quarter$ = this.route.queryParams.pipe(
      map(params => params['quarter'] ? params['quarter'] : null),
      distinctUntilChanged()
    );

    team$.pipe(
      combineLatestWith(quarter$, this.reload$),
      debounceTime(1),
      switchMap(([teamId, quarter, _]) => {
        const query: Record<string, any> = {}
        if (teamId) {
          query['team'] = teamId
        }
        if (quarter) {
          query['quarter'] = quarter
        }

        return this.sprintService.list(query).pipe(
          catchErrorMessages(this.toast),
        )
      }),
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

    this.form.valueChanges.pipe(
      distinctUntilChanged(),
      takeUntilDestroyed()
    ).subscribe(data => {
      this.router.navigate([], {queryParams: data, queryParamsHandling: 'merge'})
    })

    this.route.queryParams.pipe(
      distinctUntilChanged(),
      takeUntilDestroyed()
    ).subscribe(data => {
      this.form.patchValue(data)
    })

    this.sprintService.quarters().pipe(
      catchErrorMessages(this.toast),
      takeUntilDestroyed()
    ).subscribe(data => {
      this.allQuarters = data
      this.quartersValues = this.allQuarters.map(q => ({id: q.id, title: q.title}))
    })
  }

  add() {
    this.sprintModal.show()
      .then(sprint => {
        if (sprint) {
          this.reload()
        }
      })
      .catch(() => {})
  }

  reload() {
    this.reload$.next(null)
  }

  onDeleteSprint(sprint: Sprint) {
    this.sprints = this.sprints.filter(s => s.id !== sprint.id)
  }
}
