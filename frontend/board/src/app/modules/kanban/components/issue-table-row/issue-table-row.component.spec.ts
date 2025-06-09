import { ComponentFixture, TestBed } from '@angular/core/testing';

import { IssueTableRowComponent } from './issue-table-row.component';

describe('IssueTableRowComponent', () => {
  let component: IssueTableRowComponent;
  let fixture: ComponentFixture<IssueTableRowComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [IssueTableRowComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(IssueTableRowComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
