import { SelectValue } from "../../ui/models/select-value";
import {KanbanLabel} from "./kanban-label";
import { User } from "./user";

export enum BindStatus {
  BACKLOG = 'backlog',
  IN_PROGRESS = 'in_progress',
  DONE = 'done'
}

export const BIND_STATUS_LABELS: Record<BindStatus, string> = {
  [BindStatus.BACKLOG]: 'В плане',
  [BindStatus.IN_PROGRESS]: 'В прогрессе',
  [BindStatus.DONE]: 'Готово',
}

export const BIND_STATUS_VALUES: SelectValue[] = [
  {id: BindStatus.BACKLOG, title: BIND_STATUS_LABELS[BindStatus.BACKLOG]},
  {id: BindStatus.IN_PROGRESS, title: BIND_STATUS_LABELS[BindStatus.IN_PROGRESS]},
  {id: BindStatus.DONE, title: BIND_STATUS_LABELS[BindStatus.DONE]},
]

export type KanbanIssue = {
  id: string
  iid: string
  title: string
  type: string
  webUrl: string
  assignees: User[],
  labels: KanbanLabel[]
  projectId: number
  projectName: string
  milestone: {
    id: number
    title: string
    webPath: string
  }
  taskType: KanbanLabel | null
  estimateDev: number | null
  estimateQA: number | null
  bindingId: number | null
  bindStatus: BindStatus | null
}
