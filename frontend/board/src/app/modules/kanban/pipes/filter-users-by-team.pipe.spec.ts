import { FilterUsersByTeamPipe } from './filter-users-by-team.pipe';

describe('FilterUsersByTeamPipe', () => {
  it('create an instance', () => {
    const pipe = new FilterUsersByTeamPipe();
    expect(pipe).toBeTruthy();
  });
});
