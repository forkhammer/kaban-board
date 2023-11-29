import { FilterUsersByTextPipe } from './filter-users-by-text.pipe';

describe('FilterUsersByTextPipe', () => {
  it('create an instance', () => {
    const pipe = new FilterUsersByTextPipe();
    expect(pipe).toBeTruthy();
  });
});
