import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SprintCardComponent } from './sprint-card.component';

describe('SprintCardComponent', () => {
  let component: SprintCardComponent;
  let fixture: ComponentFixture<SprintCardComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SprintCardComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SprintCardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
