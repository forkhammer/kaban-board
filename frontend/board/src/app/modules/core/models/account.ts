export interface Account {
  id: number;
  name: string;
  username: string;
  is_active: boolean;
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
