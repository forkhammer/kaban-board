import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SprintModalComponent } from './sprint-modal.component';

describe('SprintModalComponent', () => {
  let component: SprintModalComponent;
  let fixture: ComponentFixture<SprintModalComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SprintModalComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SprintModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
