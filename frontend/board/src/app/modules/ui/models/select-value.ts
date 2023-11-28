export type SelectValueIdentity = string | number | boolean;

export interface SelectValue {
  id: SelectValueIdentity;
  title: string;
  [key: string]: any
}
