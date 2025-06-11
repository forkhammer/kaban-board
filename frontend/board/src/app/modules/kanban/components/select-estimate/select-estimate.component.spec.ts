import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SelectEstimateComponent } from './select-estimate.component';

describe('SelectEstimateComponent', () => {
  let component: SelectEstimateComponent;
  let fixture: ComponentFixture<SelectEstimateComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SelectEstimateComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SelectEstimateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
