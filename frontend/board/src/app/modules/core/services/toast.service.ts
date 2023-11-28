import { Injectable } from '@angular/core';
import { Toast } from "../models/toast";

@Injectable({
  providedIn: 'root'
})
export class ToastService {
  toasts: Toast[] = [];

  constructor() { }

  /**
   * Создает уведомление
   * @param text
   * @param className
   * @param delay
   */
  createToast(text: string, className?: string, delay: number = 5000) {
    return {
      text,
      delay,
      className,
    }
  }

  createSuccessToast(text: string, delay: number = 5000) {
    return this.createToast(text, 'success', delay);
  }

  createErrorToast(text: string, delay: number = 5000) {
    return this.createToast(text, 'error', delay);
  }

  createWarningToast(text: string, delay: number = 5000) {
    return this.createToast(text, 'warning', delay);
  }

  /**
   * Показывает уведомление
   * @param toast
   */
  show(toast: Toast) {
    this.toasts.push(toast);
  }

  remove(toast: Toast) {
    this.toasts = this.toasts.filter(t => t !== toast);
  }

  clear() {
    this.toasts.splice(0, this.toasts.length);
  }

  /**
   * Показывает несколько сообщений
   * @param messages
   */
  showMessages(messages: string[], className?: string) {
     messages.forEach(msg => {
       this.show(this.createToast(msg, className));
     });
  }
}
