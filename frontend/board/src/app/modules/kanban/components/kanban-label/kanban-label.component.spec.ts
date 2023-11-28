import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KanbanLabelComponent } from './kanban-label.component';

describe('KanbanLabelComponent', () => {
  let component: KanbanLabelComponent;
  let fixture: ComponentFixture<KanbanLabelComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [KanbanLabelComponent]
    });
    fixture = TestBed.createComponent(KanbanLabelComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
