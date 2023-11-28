import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AdminUserListComponent } from './admin-user-list.component';

describe('AdminUserListComponent', () => {
  let component: AdminUserListComponent;
  let fixture: ComponentFixture<AdminUserListComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [AdminUserListComponent]
    });
    fixture = TestBed.createComponent(AdminUserListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
