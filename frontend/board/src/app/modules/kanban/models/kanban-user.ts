import {KanbanIssue} from "./kanban-issue";

export type KanbanUser = {
  id: number
  name: string
  username: string
  avatarUrl: string
  issues: KanbanIssue[]
  teams: number[]
}

export type KanbanUserResponse = {
  users: KanbanUser[],
  updateTime: string | null
}
