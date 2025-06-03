import { Quarter } from "./quarter"
import { Team } from "./team"

export enum SprintStatus {
  RUNNING = 'running',
  COMPLETED = 'completed',
  WAITING = 'waiting'
}

export type Sprint = {
  id: number
  title: string
  start_date: string
  end_date: string
  team_id: number
  team: Team
  status: SprintStatus
  quarter: Quarter
}
