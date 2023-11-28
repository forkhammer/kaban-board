import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KanbanColumnModalComponent } from './kanban-column-modal.component';

describe('KanbanColumnModalComponent', () => {
  let component: KanbanColumnModalComponent;
  let fixture: ComponentFixture<KanbanColumnModalComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [KanbanColumnModalComponent]
    });
    fixture = TestBed.createComponent(KanbanColumnModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
