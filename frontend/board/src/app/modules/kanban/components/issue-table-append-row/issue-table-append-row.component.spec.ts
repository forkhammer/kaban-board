import { ComponentFixture, TestBed } from '@angular/core/testing';

import { IssueTableAppendRowComponent } from './issue-table-append-row.component';

describe('IssueTableAppendRowComponent', () => {
  let component: IssueTableAppendRowComponent;
  let fixture: ComponentFixture<IssueTableAppendRowComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [IssueTableAppendRowComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(IssueTableAppendRowComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
