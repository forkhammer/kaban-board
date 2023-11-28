import { ComponentFixture, TestBed } from '@angular/core/testing';

import { HeaderAccountComponent } from './header-account.component';

describe('HeaderAccountComponent', () => {
  let component: HeaderAccountComponent;
  let fixture: ComponentFixture<HeaderAccountComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [HeaderAccountComponent]
    });
    fixture = TestBed.createComponent(HeaderAccountComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
