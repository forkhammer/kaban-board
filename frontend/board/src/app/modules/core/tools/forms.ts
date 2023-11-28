import { FormGroup } from '@angular/forms';

/**
 * Отмечает поля формы в статус Dirty если они были изменены
 * Это нужно делать потому что form.patchValue не делает это самостоятельно и
 * uiInputWrap начинает работать некорректно
 */
export function markControlIfDirty(form: FormGroup) {
  Object.keys(form.controls).forEach(key => {
    const control = form.get(key);
    if (control instanceof FormGroup) {
      markControlIfDirty(control);
    } else {
      if (control?.value !== null && control?.value !== undefined) {
        control.markAsDirty();
      }
    }
  });
}
