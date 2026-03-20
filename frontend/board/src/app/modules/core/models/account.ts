export enum AccountRole {
  ADMIN = 'admin',
  EMPLOYEE = 'employee',
  VIEWER = 'viewer',
}

export interface Account {
  id: number;
  name: string;
  username: string;
  isActive: boolean;
  avatarUrl: string | null
  authProvider: string
  role: AccountRole
}

export interface AccountAuthResult {
  result?: boolean;
  message?: string;
  user: Account | null;
  errors?: string[];
}

export interface OnlineAccount {
  id: number;
  name: string;
  avatarUrl: string | null;
}

export interface RegistrationRequest {
  username: string;
  password: string;
  first_name?: string;
  last_name?: string;
}

export interface RegistrationResult {
  user: Account;
  message: string;
}
