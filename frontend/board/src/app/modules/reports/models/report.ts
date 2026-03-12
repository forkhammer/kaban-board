export interface BurndownDataPoint {
  date: string;
  total_scope: number;
  remaining_dev: number;
}

export interface BurndownReport {
  sprint_title: string;
  start_date: string;
  end_date: string;
  data_points: BurndownDataPoint[];
}
