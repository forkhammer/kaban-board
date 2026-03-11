import { Group } from "./group"
import { SelectValue } from "../../ui/models/select-value";


export enum AccountRole {
  ADMIN = 'admin',
  EMPLOYEE = 'employee',
  VIEWER = 'viewer'
}

export const ACCOUNT_ROLE_LABELS: Record<AccountRole, string> = {
  [AccountRole.ADMIN]: 'Админ',
  [AccountRole.EMPLOYEE]: 'Сотрудник',
  [AccountRole.VIEWER]: 'Наблюдатель',
}

export const ACCOUNT_ROLE_VALUES: SelectValue[] = [
  {id: AccountRole.ADMIN, title: ACCOUNT_ROLE_LABELS[AccountRole.ADMIN]},
  {id: AccountRole.EMPLOYEE, title: ACCOUNT_ROLE_LABELS[AccountRole.EMPLOYEE]},
  {id: AccountRole.VIEWER, title: ACCOUNT_ROLE_LABELS[AccountRole.VIEWER]},
]

export type User = {
  id: number
  name: string
  username: string
  avatar_url: string
  is_visible: boolean
  groups: Group[]
  account: {
    id: number
    name: string
    role: AccountRole
  } | null
}

