import {KanbanLabel} from "./kanban-label";

export type KanbanIssue = {
  id: string
  iid: string
  title: string
  type: string
  webUrl: string
  labels: {
    nodes: KanbanLabel[]
  }
  projectId: number
  projectName: string
  milestone: {
    id: number
    title: string
    webPath: string
  }
}
