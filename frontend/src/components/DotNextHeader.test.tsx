import { describe, expect, it } from "vitest";

import { navSearchOptions } from "./DotNextHeader";

// Mock translator: nav labels are keys ("menu.users") -> localized text ("Users").
const t = (k: string): string =>
  (({ "menu.dashboard": "Dashboard", "menu.users": "Users", "menu.domains": "Domains" }) as Record<string, string>)[
    k
  ] ?? k;

const nav = [
  { label: "menu.dashboard", path: "/admin/dashboard" },
  { label: "menu.users", path: "/admin/users" },
  { label: "menu.domains", path: "/admin/domains" },
];

describe("navSearchOptions (GH #455)", () => {
  it("displays the TRANSLATED label, not the raw key", () => {
    const opts = navSearchOptions(nav, "", t);
    expect(opts.map((o) => o.label)).toEqual(["Dashboard", "Users", "Domains"]);
    expect(opts.map((o) => o.value)).toEqual([
      "page:/admin/dashboard",
      "page:/admin/users",
      "page:/admin/domains",
    ]);
  });

  it("filters on the TRANSLATED text", () => {
    expect(navSearchOptions(nav, "user", t).map((o) => o.label)).toEqual(["Users"]);
  });

  it("does NOT match against the raw key (the bug)", () => {
    // "menu." prefixes every key; filtering on the localized text must ignore it.
    expect(navSearchOptions(nav, "menu", t)).toHaveLength(0);
  });
});
