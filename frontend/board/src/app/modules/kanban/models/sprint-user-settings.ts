export type SprintUserSettingsUser = {
  id: number
  name: string
  avatar_url: string
}

export type SprintUserSettings = {
  id: number
  sprint_id: number
  user_id: number
  user: SprintUserSettingsUser
  hours_per_user: number
}

export type SaveSprintUserSettingsRequest = {
  user_id: number
  hours_per_user: number
}
