import { Group } from "./group"

export type User = {
  id: number
  name: string
  username: string
  avatar_url: string
  is_visible: boolean
  groups: Group[]
}

