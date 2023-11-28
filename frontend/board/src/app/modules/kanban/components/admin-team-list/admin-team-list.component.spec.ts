import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AdminTeamListComponent } from './admin-team-list.component';

describe('AdminTeamListComponent', () => {
  let component: AdminTeamListComponent;
  let fixture: ComponentFixture<AdminTeamListComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [AdminTeamListComponent]
    });
    fixture = TestBed.createComponent(AdminTeamListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
