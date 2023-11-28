import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {IndexPageComponent} from "./components/index-page/index-page.component";
import {AuthPageComponent} from "./components/auth-page/auth-page.component";

const routes: Routes = [
  {path:'', component: IndexPageComponent, pathMatch: 'full'},
  {path:'auth', component: AuthPageComponent, pathMatch: 'full'},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
