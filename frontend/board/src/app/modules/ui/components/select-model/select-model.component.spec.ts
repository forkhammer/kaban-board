import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SelectModelComponent } from './select-model.component';

describe('SelectModelComponent', () => {
  let component: SelectModelComponent;
  let fixture: ComponentFixture<SelectModelComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SelectModelComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SelectModelComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
