export type GitlabAuthResponse = {
  enabled: boolean
  url: string | null
}

export type GitlabAuthCallbackResponse = {
  token: string
  user: {
    id: number
    username: string
    name: string
    isActive: boolean
    avatarUrl: string | null
    authProvider: string
  }
}
