export interface BurndownDataPoint {
  date: string;
  total_scope: number;
  remaining_dev: number;
  total_scope_qa: number;
  remaining_qa: number;
}

export interface BurndownReport {
  sprint_title: string;
  start_date: string;
  end_date: string;
  data_points: BurndownDataPoint[];
}

export interface BurnupDataPoint {
  date: string;
  scope_dev: number;
  completed_dev: number;
  scope_qa: number;
  completed_qa: number;
}

export interface BurnupReport {
  sprint_title: string;
  start_date: string;
  end_date: string;
  data_points: BurnupDataPoint[];
}

export interface WipDataPoint {
  label: string;
  wip_count: number;
}

export interface WipReport {
  data_points: WipDataPoint[];
}

export interface SprintStats {
  capacity: number;
  planned_dev: number;
  planned_qa: number;
  velocity_dev: number;
  velocity_qa: number;
}
