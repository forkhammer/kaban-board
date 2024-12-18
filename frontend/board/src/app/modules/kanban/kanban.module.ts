import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { UserBoardComponent } from './components/user-board/user-board.component';
import { KanbanBoardComponent } from './components/kanban-board/kanban-board.component';
import {UiModule} from "../ui/ui.module";
import { KanbanIssueComponent } from './components/kanban-issue/kanban-issue.component';
import { KanbanLabelComponent } from './components/kanban-label/kanban-label.component';
import { FilterByLabelsPipe } from './pipes/filter-by-labels.pipe';
import {FilterUnusedByLabelsPipe} from "./pipes/filter-unused-by-labels.pipe";
import { UserListComponent } from './components/user-list/user-list.component';
import { UserCardComponent } from './components/user-card/user-card.component';
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import {FontAwesomeModule} from "@fortawesome/angular-fontawesome";
import { SearchUsersPipe } from './pipes/search-users.pipe';
import { FilterUsersByTeamPipe } from './pipes/filter-users-by-team.pipe';
import { KanbanColumnComponent } from './components/kanban-column/kanban-column.component';
import {RouterLink} from "@angular/router";
import { KanbanColumnModalComponent } from './components/kanban-column-modal/kanban-column-modal.component';
import { FilterColumnsByTeamPipe } from './pipes/filter-columns-by-team.pipe';
import { OffcanvasAdminComponent } from './components/offcanvas-admin/offcanvas-admin.component';
import {BootstrapUiModule} from "../bootstrap-ui/bootstrap-ui.module";
import { AdminTeamListComponent } from './components/admin-team-list/admin-team-list.component';
import { AdminTeamCardComponent } from './components/admin-team-card/admin-team-card.component';
import { AdminUserListComponent } from './components/admin-user-list/admin-user-list.component';
import { AdminUserCardComponent } from './components/admin-user-card/admin-user-card.component';
import { AdminProjectListComponent } from './components/admin-project-list/admin-project-list.component';
import { AdminProjectCardComponent } from './components/admin-project-card/admin-project-card.component';
import { FilterUsersByTextPipe } from './pipes/filter-users-by-text.pipe';
import { FilterIssuesByTextPipe } from './pipes/filter-issues-by-text.pipe';
import { FilterUsersByNamePipe } from './pipes/filter-users-by-name.pipe';
import { FilterProjectsByNamePipe } from './pipes/filter-project-by-name.pipe';
import { DragDropModule } from '@angular/cdk/drag-drop';
import { AdminSettingsComponent } from './components/admin-settings/admin-settings.component';
import { AdminLabelListComponent } from './components/admin-label-list/admin-label-list.component'
import { FilterLabelsByNamePipe } from './pipes/filter-labels-by-name.pipe';
import { AdminLabelCardComponent } from './components/admin-label-card/admin-label-card.component';
import { UserGroupCardComponent } from './components/user-group-card/user-group-card.component';
import { FilterUsersByGroupPipe } from './pipes/filter-users-by-group.pipe';


@NgModule({
  declarations: [
    UserBoardComponent,
    KanbanBoardComponent,
    KanbanIssueComponent,
    KanbanLabelComponent,
    FilterByLabelsPipe,
    FilterUnusedByLabelsPipe,
    UserListComponent,
    UserCardComponent,
    SearchUsersPipe,
    FilterUsersByTeamPipe,
    KanbanColumnComponent,
    KanbanColumnModalComponent,
    FilterColumnsByTeamPipe,
    OffcanvasAdminComponent,
    AdminTeamListComponent,
    AdminTeamCardComponent,
    AdminUserListComponent,
    AdminUserCardComponent,
    AdminProjectListComponent,
    AdminProjectCardComponent,
    FilterUsersByTextPipe,
    FilterIssuesByTextPipe,
    FilterUsersByNamePipe,
    FilterProjectsByNamePipe,
    AdminSettingsComponent,
    AdminLabelListComponent,
    FilterLabelsByNamePipe,
    AdminLabelCardComponent,
    UserGroupCardComponent,
    FilterUsersByGroupPipe,
  ],
  exports: [
    UserBoardComponent,
    KanbanBoardComponent,
    OffcanvasAdminComponent,
  ],
  imports: [
    CommonModule,
    UiModule,
    BootstrapUiModule,
    FormsModule,
    ReactiveFormsModule,
    FontAwesomeModule,
    RouterLink,
    DragDropModule,
  ]
})
export class KanbanModule { }
