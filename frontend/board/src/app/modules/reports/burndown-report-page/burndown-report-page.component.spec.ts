import { ComponentFixture, TestBed } from '@angular/core/testing';

import { BurndownReportPageComponent } from './burndown-report-page.component';

describe('BurndownReportPageComponent', () => {
  let component: BurndownReportPageComponent;
  let fixture: ComponentFixture<BurndownReportPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [BurndownReportPageComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(BurndownReportPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
