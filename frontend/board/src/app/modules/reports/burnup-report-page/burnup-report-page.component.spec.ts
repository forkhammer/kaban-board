import { ComponentFixture, TestBed } from '@angular/core/testing';

import { BurnupReportPageComponent } from './burnup-report-page.component';

describe('BurnupReportPageComponent', () => {
  let component: BurnupReportPageComponent;
  let fixture: ComponentFixture<BurnupReportPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [BurnupReportPageComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(BurnupReportPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
