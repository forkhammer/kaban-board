import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AuthGitlabCallbackPageComponent } from './auth-gitlab-callback-page.component';

describe('AuthGitlabCallbackPageComponent', () => {
  let component: AuthGitlabCallbackPageComponent;
  let fixture: ComponentFixture<AuthGitlabCallbackPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AuthGitlabCallbackPageComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(AuthGitlabCallbackPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
