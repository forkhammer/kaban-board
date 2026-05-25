import { Component, inject } from "@angular/core";
import { ActivatedRoute, Router } from "@angular/router";
import { distinctUntilChanged, map } from "rxjs";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";

@Component({
  selector: "app-search-page",
  standalone: false,
  templateUrl: "./search-page.component.html",
  styleUrl: "./search-page.component.scss",
})
export class SearchPageComponent {
  private route = inject(ActivatedRoute);

  searchQuery: string | null = null;

  constructor() {
    this.route.queryParams
      .pipe(
        map(
          (params) =>
            (params["search"] ? String(params["search"]) : null) as
              | string
              | null,
        ),
        distinctUntilChanged(),
        takeUntilDestroyed(),
      )
      .subscribe((search) => {
        this.searchQuery = search;
      });
  }
}
