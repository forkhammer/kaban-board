import { ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';

import { UserCardComponent } from './user-card.component';
import { UserWorkload } from '../../../reports/models/report';

describe('UserCardComponent', () => {
  let component: UserCardComponent;
  let fixture: ComponentFixture<UserCardComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [UserCardComponent]
    });
    fixture = TestBed.createComponent(UserCardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should not render workload element when workload is undefined', () => {
    expect(component.workload).toBeUndefined();
    expect(fixture.debugElement.query(By.css('.workload'))).toBeNull();
  });

  it('should compute loadPercent 50 and success color for 20/40', () => {
    component.workload = { user_id: 1, capacity: 40, planned: 20 } as UserWorkload;
    expect(component.loadPercent).toBe(50);
    expect(component.loadColorClass).toBe('success');
  });

  it('should compute loadPercent 90 and warning color for 36/40', () => {
    component.workload = { user_id: 1, capacity: 40, planned: 36 } as UserWorkload;
    expect(component.loadPercent).toBe(90);
    expect(component.loadColorClass).toBe('warning');
  });

  it('should compute loadPercent 125 and danger color for 50/40', () => {
    component.workload = { user_id: 1, capacity: 40, planned: 50 } as UserWorkload;
    expect(component.loadPercent).toBe(125);
    expect(component.loadColorClass).toBe('danger');
  });

  it('should compute loadPercent 0 (not NaN) for 0/0', () => {
    component.workload = { user_id: 1, capacity: 0, planned: 0 } as UserWorkload;
    expect(component.loadPercent).toBe(0);
    expect(component.loadColorClass).toBe('success');
  });

  it('should render workload element after setting workload', () => {
    component.workload = { user_id: 1, capacity: 40, planned: 20 } as UserWorkload;
    fixture.detectChanges();
    const el = fixture.debugElement.query(By.css('.workload'));
    expect(el).toBeTruthy();
    expect(fixture.nativeElement.textContent).toContain('50%');
  });
});
