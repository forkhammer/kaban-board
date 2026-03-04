import { SelectValue } from "../../ui/models/select-value";
import { Epic } from "./epic";
import {KanbanLabel} from "./kanban-label";
import { Release } from "./release";
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

export enum IssuePriority {
  LOWEST = 'lowest',
  LOW = 'low',
  MEDIUM = 'medium',
  HIGH = 'high',
  CRITICAL = 'critical'
}

export const ISSUE_PRIORITY_LABELS: Record<IssuePriority, string> = {
  [IssuePriority.LOWEST]: 'Самый низкий',
  [IssuePriority.LOW]: 'Низкий',
  [IssuePriority.MEDIUM]: 'Средний',
  [IssuePriority.HIGH]: 'Высокий',
  [IssuePriority.CRITICAL]: 'Критический'
}

export const ISSUE_PRIORITY_VALUES: SelectValue[] = [
  {id: IssuePriority.LOWEST, title: ISSUE_PRIORITY_LABELS[IssuePriority.LOWEST]},
  {id: IssuePriority.LOW, title: ISSUE_PRIORITY_LABELS[IssuePriority.LOW]},
  {id: IssuePriority.MEDIUM, title: ISSUE_PRIORITY_LABELS[IssuePriority.MEDIUM]},
  {id: IssuePriority.HIGH, title: ISSUE_PRIORITY_LABELS[IssuePriority.HIGH]},
  {id: IssuePriority.CRITICAL, title: ISSUE_PRIORITY_LABELS[IssuePriority.CRITICAL]},
]

export type KanbanIssue = {
  id: string
  iid: string
  title: string
  type: string
  webUrl: string
  assignees: User[],
  assignee: User | null,
  labels: KanbanLabel[]
  projectId: number
  projectName: string
  release: Release | null
  epic: Epic | null
  taskType: KanbanLabel | null
  estimateDev: number | null
  estimateQA: number | null
  bindingId: number | null
  bindStatus: BindStatus | null
  priority: IssuePriority | null,
  comment: string | null
}

export type IssueBinding = {
  id: string
  iid: string
  title: string
  type: string
  webUrl: string
  assignees: User[],
  assignee: User | null,
  labels: KanbanLabel[]
  projectId: number
  projectName: string
  release: Release | null
  epic: Epic | null
  taskType: KanbanLabel | null
  estimateDev: number | null
  estimateQA: number | null
  bindingId: number
  bindStatus: BindStatus
  priority: IssuePriority | null,
  comment: string | null
}
