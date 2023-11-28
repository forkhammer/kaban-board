import { ComponentFixture, TestBed } from '@angular/core/testing';

import { UserBoardComponent } from './user-board.component';

describe('UserBoardComponent', () => {
  let component: UserBoardComponent;
  let fixture: ComponentFixture<UserBoardComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [UserBoardComponent]
    });
    fixture = TestBed.createComponent(UserBoardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
