import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AdminLabelCardComponent } from './admin-label-card.component';

describe('AdminLabelCardComponent', () => {
  let component: AdminLabelCardComponent;
  let fixture: ComponentFixture<AdminLabelCardComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [AdminLabelCardComponent]
    });
    fixture = TestBed.createComponent(AdminLabelCardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
