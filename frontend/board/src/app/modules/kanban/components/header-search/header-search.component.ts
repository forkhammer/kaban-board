import { Component, inject } from "@angular/core";
import { Router } from "@angular/router";

@Component({
  selector: "app-header-search",
  standalone: false,
  templateUrl: "./header-search.component.html",
  styleUrl: "./header-search.component.scss",
})
export class HeaderSearchComponent {
  private router = inject(Router);

  searchValue = "";

  onSearch(event: Event) {
    event.preventDefault();
    const value = this.searchValue.trim();
    this.router.navigate(["/search"], { queryParams: { search: value } });
  }

  onKeydown(event: KeyboardEvent) {
    if (event.key === "Enter") {
      this.onSearch(event);
    }
  }

  clear() {
    this.searchValue = "";
    this.onSearch(new Event("search"));
  }
}
