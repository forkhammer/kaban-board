import { ComponentFixture, TestBed } from '@angular/core/testing';

import { InputAlertComponent } from './input-alert.component';

describe('InputAlertComponent', () => {
  let component: InputAlertComponent;
  let fixture: ComponentFixture<InputAlertComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ InputAlertComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(InputAlertComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
