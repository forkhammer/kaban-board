import {KanbanLabel} from "./kanban-label";
import { User } from "./user";

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
}
