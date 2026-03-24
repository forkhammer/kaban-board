import { Component, inject } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';
import { CreateIssueBindingRequest, IssueBindingService } from '../../services/issue-binding.service';
import { ProjectService } from '../../services/project.service';
import { IssueBindingModalData } from '../../services/issue-binding-modal.service';
import { catchError, finalize } from 'rxjs/operators';
import { of } from 'rxjs';
import { UserService } from '../../services/user.service';
import { ToastService } from 'src/app/modules/core/services/toast.service';

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
  protected toast = inject(ToastService)

  form: FormGroup;
  isLoading = false;
  errors: string[] = [];

  constructor() {
    this.form = this.fb.group({
      id: [null],
      title: ['', [Validators.required]],
      project: [null, [Validators.required]],
      assignee: [null],
      sprint: [null],
      createInTracker: [false],
    });
  }

  /**
   * Инициализация модального окна с начальными данными
   */
  init(data: IssueBindingModalData): void {
    this.form.patchValue({
      id: data.id,
      title: data.title,
      assignee: data.assigneeId,
      sprint: data.sprintId,
      project: data.projectId
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

    const data: CreateIssueBindingRequest & { id: number } = {
      id: this.form.value.id,
      title: this.form.value.title,
      project: this.form.value.project,
      assignee: this.form.value.assignee,
      sprint: this.form.value.sprint,
      createInTracker: this.form.value.createInTracker,
    };

    if (data.id) {
      this.issueBindingService.save(data)
        .pipe(
          finalize(() => this.isLoading = false),
          catchError((error) => {
            if (error.error?.message) {
              this.errors = [error.error.message];
            } else if (error.error?.error) {
              this.errors = [error.error.error];
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
    } else {
      this.issueBindingService.create(data)
        .pipe(
          finalize(() => this.isLoading = false),
          catchError((error) => {
            if (error.error?.message) {
              this.errors = [error.error.message];
            } else if (error.error?.error) {
              this.errors = [error.error.error];
            } else {
              this.errors = ['An error occurred while creating the issue binding'];
            }

            this.toast.showMessages(this.errors, 'error')
            return of(null);
          })
        )
        .subscribe((result) => {
          if (result) {
            this.modal.close(result);
          }
        });
    }
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
