import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SelectModelMultipleComponent } from './select-model-multiple.component';

describe('SelectModelComponent', () => {
  let component: SelectModelMultipleComponent;
  let fixture: ComponentFixture<SelectModelMultipleComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SelectModelMultipleComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SelectModelMultipleComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
