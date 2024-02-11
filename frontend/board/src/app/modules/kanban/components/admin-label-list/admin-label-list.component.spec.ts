import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AdminLabelListComponent } from './admin-label-list.component';

describe('AdminLabelListComponent', () => {
  let component: AdminLabelListComponent;
  let fixture: ComponentFixture<AdminLabelListComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [AdminLabelListComponent]
    });
    fixture = TestBed.createComponent(AdminLabelListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
