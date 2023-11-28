import {Component, OnDestroy, OnInit} from '@angular/core';
import {finalize, Subject} from "rxjs";
import {Project} from "../../models/project";
import {ProjectService} from "../../services/project.service";
import {takeUntil} from "rxjs/operators";

@Component({
  selector: 'app-admin-project-list',
  templateUrl: './admin-project-list.component.html',
  styleUrls: ['./admin-project-list.component.scss']
})
export class AdminProjectListComponent implements OnInit, OnDestroy {
  private destroy$ = new Subject()
  public projects: Project[] = []
  public isLoading = true

  constructor(
    private projectService: ProjectService
  ) {
  }

  ngOnInit() {
    this.projectService.all().pipe(
      finalize(() => this.isLoading = false),
      takeUntil(this.destroy$)
    ).subscribe(data => {
      this.projects = data as Project[]
    })
  }

  ngOnDestroy() {
    this.destroy$.next(null)
    this.destroy$.complete()
  }

  trackByProject(_: number, project: Project): number {
    return project.id
  }
}
