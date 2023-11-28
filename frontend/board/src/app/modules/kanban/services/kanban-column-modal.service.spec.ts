import { TestBed } from '@angular/core/testing';

import { KanbanColumnModalService } from './kanban-column-modal.service';

describe('KanbanColumnModalService', () => {
  let service: KanbanColumnModalService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(KanbanColumnModalService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
