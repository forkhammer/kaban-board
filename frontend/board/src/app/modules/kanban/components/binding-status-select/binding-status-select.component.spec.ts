import { ComponentFixture, TestBed } from '@angular/core/testing';

import { BindingStatusSelectComponent } from './binding-status-select.component';

describe('BindingStatusSelectComponent', () => {
  let component: BindingStatusSelectComponent;
  let fixture: ComponentFixture<BindingStatusSelectComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [BindingStatusSelectComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(BindingStatusSelectComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
