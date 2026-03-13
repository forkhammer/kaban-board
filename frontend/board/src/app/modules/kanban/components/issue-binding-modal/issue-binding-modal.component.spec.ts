import { ComponentFixture, TestBed } from '@angular/core/testing';

import { IssueBindingModalComponent } from './issue-binding-modal.component';

describe('IssueBindingModalComponent', () => {
  let component: IssueBindingModalComponent;
  let fixture: ComponentFixture<IssueBindingModalComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [IssueBindingModalComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(IssueBindingModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
