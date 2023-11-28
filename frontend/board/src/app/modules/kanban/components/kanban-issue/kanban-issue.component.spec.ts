import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KanbanIssueComponent } from './kanban-issue.component';

describe('KanbanIssueComponent', () => {
  let component: KanbanIssueComponent;
  let fixture: ComponentFixture<KanbanIssueComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [KanbanIssueComponent]
    });
    fixture = TestBed.createComponent(KanbanIssueComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
