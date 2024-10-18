import { ComponentFixture, TestBed } from '@angular/core/testing';

import { OtherGroupCardComponent } from './other-group-card.component';

describe('OtherGroupCardComponent', () => {
  let component: OtherGroupCardComponent;
  let fixture: ComponentFixture<OtherGroupCardComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [OtherGroupCardComponent]
    });
    fixture = TestBed.createComponent(OtherGroupCardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
