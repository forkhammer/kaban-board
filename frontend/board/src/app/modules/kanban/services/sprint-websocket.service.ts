import { inject, Injectable } from '@angular/core';
import { Subject, Subscription, retry } from 'rxjs';
import { webSocket, WebSocketSubject } from 'rxjs/webSocket';
import { JWTService } from '../../core/services/jwt.service';
import { CoreConfigService } from '../../core/config';

export type WSEventType = 'binding_updated' | 'binding_created' | 'binding_deleted' | 'binding_ordering';

export interface BindingDeletedEventData {
  id: number
  sprint_id: number
}

export interface BindingWSEvent {
  type: WSEventType;
  data: any;
}

@Injectable({ providedIn: 'root' })
export class SprintWebsocketService {
  private config = inject(CoreConfigService);
  private jwt = inject(JWTService);

  private socket$: WebSocketSubject<BindingWSEvent> | null = null;
  private subscription: Subscription | null = null;
  private eventsSubject = new Subject<BindingWSEvent>();
  readonly events$ = this.eventsSubject.asObservable();

  connect(sprintId: number): void {
    this.disconnect();
    const wsUrl = this.config.apiUrl.replace(/^http/, 'ws') + `/sprint/${sprintId}/ws?token=${this.jwt.access}`;
    this.socket$ = webSocket<BindingWSEvent>(wsUrl);
    this.subscription = this.socket$.pipe(
      retry({ delay: 3000 })
    ).subscribe({
      next: event => this.eventsSubject.next(event),
      error: () => {},
    });
  }

  disconnect(): void {
    this.subscription?.unsubscribe();
    this.subscription = null;
    this.socket$?.complete();
    this.socket$ = null;
  }
}
