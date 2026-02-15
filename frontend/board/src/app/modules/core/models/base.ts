export interface BaseModel {
  id: number | string;
}

export interface Pagination<T> {
  count: number;
  pages: number;
  page: number;
  start_index: number;
  end_index: number;
  results: T[];
}

export interface BaseTitleModel extends BaseModel {
  title: string;
}
