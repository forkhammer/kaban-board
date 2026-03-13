import { Component, inject } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';
import { IssueBindingService } from '../../services/issue-binding.service';
import { ProjectService } from '../../services/project.service';
import { KanbanIssue } from '../../models/kanban-issue';
import { User } from '../../models/user';
import { IssueBindingModalData } from '../../services/issue-binding-modal.service';
import { catchError, finalize } from 'rxjs/operators';
import { of } from 'rxjs';
import { UserService } from '../../services/user.service';

@Component({
  selector: 'app-issue-binding-modal',
  standalone: false,
  templateUrl: './issue-binding-modal.component.html',
  styleUrl: './issue-binding-modal.component.scss'
})
export class IssueBindingModalComponent {
  private modal = inject(NgbActiveModal);
  private fb = inject(FormBuilder);
  private issueBindingService = inject(IssueBindingService);
  protected projectService = inject(ProjectService);
  protected userService = inject(UserService)

  form: FormGroup;
  isLoading = false;
  errors: string[] = [];

  constructor() {
    this.form = this.fb.group({
      title: ['', [Validators.required]],
      project: [null, [Validators.required]],
      assignee: [null],
      sprint: [null],
    });
  }

  /**
   * Инициализация модального окна с начальными данными
   */
  init(data: IssueBindingModalData): void {
    this.form.patchValue({
      assignee: data.assigneeId,
      sprint: data.sprintId
    })
  }

  /**
   * Закрытие модального окна без сохранения
   */
  close(): void {
    this.modal.dismiss(null);
  }

  /**
   * Отправка формы и создание Issue Binding
   */
  submit(): void {
    if (this.form.invalid) {
      this.markAllAsTouched();
      return;
    }

    this.isLoading = true;
    this.errors = [];

    const data = {
      title: this.form.value.title,
      project: this.form.value.project,
      assignee: this.form.value.assignee,
      sprint: this.form.value.sprint,
    };

    this.issueBindingService.create(data)
      .pipe(
        finalize(() => this.isLoading = false),
        catchError((error) => {
          if (error.error?.message) {
            this.errors = [error.error.message];
          } else if (error.error?.errors) {
            this.errors = Object.values(error.error.errors).flat() as string[];
          } else {
            this.errors = ['An error occurred while creating the issue binding'];
          }
          return of(null);
        })
      )
      .subscribe((result) => {
        if (result) {
          this.modal.close(result);
        }
      });
  }

  /**
   * Помечает все поля формы как touched для отображения ошибок валидации
   */
  private markAllAsTouched(): void {
    Object.keys(this.form.controls).forEach(key => {
      const control = this.form.get(key);
      control?.markAsTouched();
    });
  }

  /**
   * Получение текста ошибки для поля
   */
  getErrorMessage(fieldName: string): string | null {
    const control = this.form.get(fieldName);
    if (control?.touched && control?.errors) {
      if (control.errors['required']) {
        return 'This field is required';
      }
    }
    return null;
  }
}
