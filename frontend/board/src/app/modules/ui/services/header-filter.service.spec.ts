import { TestBed } from '@angular/core/testing';

import { HeaderFilterService } from './header-filter.service';

describe('HeaderFilterService', () => {
  let service: HeaderFilterService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(HeaderFilterService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
