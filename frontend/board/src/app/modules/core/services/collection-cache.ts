import { BaseModel } from "../models/base";

export class CollectionCache<T extends BaseModel> {
  private items: Map<number | string, T> = new Map<number | string, T>();

  setItems(items: T[]) {
    items.forEach(item => {
      this.items.set(item.id, item);
    });
  }

  get(id: number | string): T | undefined {
    return this.items.get(id);
  }

  invalidate(){
    this.items = new Map<number | string, T>();
  }

  invalidateItem(id: number | string){
    this.items.delete(id);
  }
}
