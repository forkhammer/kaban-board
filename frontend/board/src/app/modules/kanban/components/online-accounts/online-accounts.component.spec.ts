import { ComponentFixture, TestBed } from '@angular/core/testing';

import { OnlineAccountsComponent } from './online-accounts.component';

describe('OnlineAccountsComponent', () => {
  let component: OnlineAccountsComponent;
  let fixture: ComponentFixture<OnlineAccountsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [OnlineAccountsComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(OnlineAccountsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
