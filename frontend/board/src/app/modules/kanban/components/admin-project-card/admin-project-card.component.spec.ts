import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AdminProjectCardComponent } from './admin-project-card.component';

describe('AdminProjectCardComponent', () => {
  let component: AdminProjectCardComponent;
  let fixture: ComponentFixture<AdminProjectCardComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [AdminProjectCardComponent]
    });
    fixture = TestBed.createComponent(AdminProjectCardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
