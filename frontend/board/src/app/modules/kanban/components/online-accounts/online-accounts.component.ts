import { Component } from '@angular/core';
import { timer } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { AccountService } from '../../../core/services/account.service';
import { OnlineAccount } from '../../../core/models/account';

@Component({
  selector: 'app-online-accounts',
  templateUrl: './online-accounts.component.html',
  styleUrl: './online-accounts.component.scss',
  standalone: false,
})
export class OnlineAccountsComponent {
  readonly maxVisible = 3;
  accounts: OnlineAccount[] = [];

  get visibleAccounts(): OnlineAccount[] {
    return this.accounts.slice(0, this.maxVisible);
  }

  get hiddenCount(): number {
    return Math.max(0, this.accounts.length - this.maxVisible);
  }

  constructor(private accountService: AccountService) {
    timer(0, 60 * 1000).pipe(
      switchMap(() => this.accountService.getOnlineAccounts()),
      takeUntilDestroyed(),
    ).subscribe(accounts => this.accounts = accounts);
  }
}
