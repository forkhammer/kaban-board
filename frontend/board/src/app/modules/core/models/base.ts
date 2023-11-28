export interface BaseModel {
  id: number | string;
}

export interface Pagination<T> {
  count?: number;
  next?: string;
  pages?: number;
  page: number;
  previous?: string;
  start_index?: number;
  end_index?: number;
  results: T[];
}

export interface BaseTitleModel extends BaseModel {
  title: string;
}
