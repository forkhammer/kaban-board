import { FilterIssuesByTextPipe } from './filter-issues-by-text.pipe';

describe('FilterIssuesByTextPipe', () => {
  it('create an instance', () => {
    const pipe = new FilterIssuesByTextPipe();
    expect(pipe).toBeTruthy();
  });
});
