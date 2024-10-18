import { Group } from "./group"

export interface Team {
  id: number
  title: string
  groups: Group[]
}
