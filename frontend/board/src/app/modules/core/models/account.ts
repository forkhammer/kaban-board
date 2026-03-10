export interface Account {
  id: number;
  name: string;
  username: string;
  isActive: boolean;
  avatarUrl: string | null
  authProvider: string
}

export interface AccountAuthResult {
  result?: boolean;
  message?: string;
  user: Account | null;
  errors?: string[];
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
