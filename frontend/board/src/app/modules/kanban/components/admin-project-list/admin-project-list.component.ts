import {Component, inject} from '@angular/core';
import {finalize} from "rxjs";
import {Project} from "../../models/project";
import {ProjectService} from "../../services/project.service";
import { FormBuilder, FormGroup } from '@angular/forms';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Component({
    selector: 'app-admin-project-list',
    templateUrl: './admin-project-list.component.html',
    styleUrls: ['./admin-project-list.component.scss'],
    standalone: false
})
export class AdminProjectListComponent {
  private projectService = inject(ProjectService)
  private fb = inject(FormBuilder)

  public projects: Project[] = []
  public isLoading = true
  public filterForm: FormGroup

  constructor(
  ) {
    this.filterForm = this.fb.group({
      search: ['']
    })

    this.projectService.all().pipe(
      finalize(() => this.isLoading = false),
      takeUntilDestroyed()
    ).subscribe(data => {
      this.projects = data as Project[]
    })
  }

  trackByProject(_: number, project: Project): number {
    return project.id
  }
}
