export type KanbanLabel = {
  id: string
  name: string
  color: string
  textColor: string
  altName: string | null
}

import { BindStatus, IssuePriority } from './kanban-issue'

export type Label = {
  id: string
  name: string
  altName: string | null
  bindingStatus: BindStatus | null
  priority: IssuePriority | null
}
