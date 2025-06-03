import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SprintListComponent } from './sprint-list.component';

describe('SprintListComponent', () => {
  let component: SprintListComponent;
  let fixture: ComponentFixture<SprintListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SprintListComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SprintListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
