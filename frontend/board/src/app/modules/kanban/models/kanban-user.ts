import { Group } from "./group";
import {KanbanIssue} from "./kanban-issue";

export type KanbanUser = {
  id: number
  name: string
  username: string
  avatarUrl: string
  issues: KanbanIssue[]
  teams: number[]
  groups: Group[]
}

export type KanbanUserResponse = {
  users: KanbanUser[],
  updateTime: string | null
}
