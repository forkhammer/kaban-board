import { ComponentFixture, TestBed } from '@angular/core/testing';

import { OffcanvasAdminComponent } from './offcanvas-admin.component';

describe('OffcanvasAdminComponent', () => {
  let component: OffcanvasAdminComponent;
  let fixture: ComponentFixture<OffcanvasAdminComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [OffcanvasAdminComponent]
    });
    fixture = TestBed.createComponent(OffcanvasAdminComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
