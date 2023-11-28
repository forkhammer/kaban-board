import { TestBed } from '@angular/core/testing';

import { GitlabSyncService } from './gitlab-sync.service';

describe('GitlabSyncService', () => {
  let service: GitlabSyncService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(GitlabSyncService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
