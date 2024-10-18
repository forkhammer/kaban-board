import { ComponentFixture, TestBed } from '@angular/core/testing';

import { UserGroupCardComponent } from './user-group-card.component';

describe('UserGroupCardComponent', () => {
  let component: UserGroupCardComponent;
  let fixture: ComponentFixture<UserGroupCardComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [UserGroupCardComponent]
    });
    fixture = TestBed.createComponent(UserGroupCardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
