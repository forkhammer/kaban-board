import { Component, DestroyRef, inject } from "@angular/core";
import { SprintService } from "../../services/sprint.service";
import { TeamService } from "../../services/team.service";
import { FormBuilder, FormGroup, Validators } from "@angular/forms";
import { ToastService } from "src/app/modules/core/services/toast.service";
import { SaveSprintRequest, Sprint } from "../../models/sprint";
import { catchErrorMessages } from "src/app/modules/core/tools/catch-error";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { NgbActiveModal } from "@ng-bootstrap/ng-bootstrap";
import { catchError, EMPTY } from "rxjs";
import { formatDate } from "@angular/common";
import { SprintUserSettingsService } from "../../services/sprint-user-settings.service";
import {
  SprintUserSettings,
  SaveSprintUserSettingsRequest,
} from "../../models/sprint-user-settings";
import { faTimes } from "@fortawesome/free-solid-svg-icons";

@Component({
  selector: "app-sprint-modal",
  standalone: false,
  templateUrl: "./sprint-modal.component.html",
  styleUrl: "./sprint-modal.component.scss",
})
export class SprintModalComponent {
  sprintService = inject(SprintService);
  fb = inject(FormBuilder);
  toast = inject(ToastService);
  destroyRef = inject(DestroyRef);
  modal = inject(NgbActiveModal);
  teamService = inject(TeamService);
  sprintUserSettingsService = inject(SprintUserSettingsService);

  readonly faTimes = faTimes;

  form: FormGroup;
  sprint: Sprint | null = null;
  isLoading = false;
  errors: any;

  userSettings: SprintUserSettings[] = [];
  selectedUserId: number | null = null;
  activeTab = 1;

  constructor() {
    this.form = this.fb.group({
      title: [""],
      start_date: [
        formatDate(this.getStartSprintDate(), "yyyy-MM-dd", "en"),
        Validators.required,
      ],
      end_date: [
        formatDate(this.getEndSprintDate(), "yyyy-MM-dd", "en"),
        Validators.required,
      ],
      hours_per_user: [60],
      team_id: [null, Validators.required],
    });
  }

  init(sprintId?: number) {
    if (sprintId) {
      this.sprintService
        .get(sprintId)
        .pipe(
          catchErrorMessages(this.toast),
          takeUntilDestroyed(this.destroyRef),
        )
        .subscribe((sprint) => {
          this.sprint = sprint;
          this.form.patchValue(sprint);
        });

      this.sprintUserSettingsService
        .getSettings(sprintId)
        .pipe(
          catchErrorMessages(this.toast),
          takeUntilDestroyed(this.destroyRef),
        )
        .subscribe((settings) => {
          this.userSettings = settings;
        });
    }
  }

  close() {
    this.modal.close(null);
  }

  submit() {
    this.isLoading = true;
    const data: SaveSprintRequest = {
      id: this.sprint?.id,
      title: this.form.value.title,
      start_date: this.form.value.start_date,
      end_date: this.form.value.end_date,
      hours_per_user: this.form.value.hours_per_user,
      team_id: this.form.value.team_id,
    };
    this.sprintService
      .save(data)
      .pipe(
        catchError((err) => {
          this.errors = err.error;
          this.isLoading = false;
          return EMPTY;
        }),
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe((sprint) => {
        this.isLoading = false;
        this.modal.close(sprint);
      });
  }

  loadUserSettings() {
    if (!this.sprint) return;
    this.sprintUserSettingsService
      .getSettings(this.sprint.id)
      .pipe(catchErrorMessages(this.toast), takeUntilDestroyed(this.destroyRef))
      .subscribe((settings) => {
        this.userSettings = settings;
      });
  }

  onHoursChange(setting: SprintUserSettings, event: Event) {
    const input = event.target as HTMLInputElement;
    const newValue = parseInt(input.value, 10);
    if (isNaN(newValue) || newValue === setting.hours_per_user) {
      input.value = String(setting.hours_per_user);
      return;
    }
    if (!this.sprint) return;

    const request: SaveSprintUserSettingsRequest = {
      user_id: setting.user_id,
      hours_per_user: newValue,
    };
    this.sprintUserSettingsService
      .saveSettings(this.sprint.id, request)
      .pipe(catchErrorMessages(this.toast), takeUntilDestroyed(this.destroyRef))
      .subscribe((updated) => {
        const idx = this.userSettings.findIndex(
          (s) => s.user_id === updated.user_id,
        );
        if (idx >= 0) {
          this.userSettings[idx] = updated;
        }
      });
  }

  addUserSetting() {
    if (!this.sprint || !this.selectedUserId) return;
    const hoursPerUser = this.sprint.hours_per_user ?? 60;
    const request: SaveSprintUserSettingsRequest = {
      user_id: this.selectedUserId,
      hours_per_user: hoursPerUser,
    };
    this.sprintUserSettingsService
      .saveSettings(this.sprint.id, request)
      .pipe(catchErrorMessages(this.toast), takeUntilDestroyed(this.destroyRef))
      .subscribe(() => {
        this.selectedUserId = null;
        this.loadUserSettings();
      });
  }

  removeUserSetting(userId: number) {
    if (!this.sprint) return;
    if (!confirm("Удалить настройку пользователя?")) return;
    this.sprintUserSettingsService
      .deleteSettings(this.sprint.id, userId)
      .pipe(catchErrorMessages(this.toast), takeUntilDestroyed(this.destroyRef))
      .subscribe(() => {
        this.loadUserSettings();
      });
  }

  getStartSprintDate() {
    return new Date();
  }

  getEndSprintDate() {
    const dt = this.getStartSprintDate();
    dt.setDate(dt.getDate() + 14);
    return dt;
  }
}
