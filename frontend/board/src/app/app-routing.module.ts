import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {IndexPageComponent} from "./components/index-page/index-page.component";
import {AuthPageComponent} from "./components/auth-page/auth-page.component";
import { SprintsPageComponent } from './components/sprints-page/sprints-page.component';
import { ReportsPageComponent } from './components/reports-page/reports-page.component';
import { AuthGitlabCallbackPageComponent } from './components/auth-gitlab-callback-page/auth-gitlab-callback-page.component';

const routes: Routes = [
  {path:'', component: IndexPageComponent, pathMatch: 'full'},
  {path:'auth/gitlab', component: AuthGitlabCallbackPageComponent, pathMatch: 'full'},
  {path:'auth', component: AuthPageComponent, pathMatch: 'full'},
  {path:'sprints', component: SprintsPageComponent, pathMatch: 'full'},
  {path:'reports', component: ReportsPageComponent, pathMatch: 'full'},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
