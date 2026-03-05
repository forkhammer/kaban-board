import { ComponentFixture, TestBed } from '@angular/core/testing';

import { BindIssueModalComponent } from './bind-issue-modal.component';

describe('BindIssueModalComponent', () => {
  let component: BindIssueModalComponent;
  let fixture: ComponentFixture<BindIssueModalComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [BindIssueModalComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(BindIssueModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
