// impersonation.ts — client side of admin act-as (ADR-0128).
// Sends the active grant ID in X-DotNext-Act-As via apiClient.
import { apiClient } from "./apiClient";
import { clearIdentity } from "./identity";

export const ACT_AS_KEY = "dotnext_act_as";
const LEGACY_ACT_AS_KEY = "dotnext_act_as";

export type ActAsState = { id: string; username: string };

export function getActAs(): ActAsState | null {
  try {
    const raw =
      sessionStorage.getItem(ACT_AS_KEY) ??
      sessionStorage.getItem(LEGACY_ACT_AS_KEY);
    return raw ? (JSON.parse(raw) as ActAsState) : null;
  } catch {
    return null;
  }
}

function setActAs(s: ActAsState | null): void {
  if (s) {
    sessionStorage.setItem(ACT_AS_KEY, JSON.stringify(s));
  } else {
    sessionStorage.removeItem(ACT_AS_KEY);
    sessionStorage.removeItem(LEGACY_ACT_AS_KEY);
  }
}

/** Start acting as a user. Returns the grant; caller should refetch auth +
 * navigate to the user shell. Runs as the real admin (no grant set yet). */
export async function startImpersonation(userId: string): Promise<ActAsState> {
  const { data } = await apiClient.post<{
    id: string;
    target_username: string;
  }>("/admin/impersonation", { user_id: userId });
  const state: ActAsState = { id: data.id, username: data.target_username };
  setActAs(state);
  // Codeberg #4: drop the cached /me so RequireAdmin re-resolves against the
  // now-effective (tenant, act-as header) identity instead of the stale admin.
  clearIdentity();
  return state;
}

/** Stop acting as a user. Clears the grant FIRST so the revoke call (and every
 * subsequent request) runs as the real admin, then best-effort revokes the
 * server-side grant. Caller should refetch auth + navigate back to admin. */
export async function stopImpersonation(): Promise<void> {
  const s = getActAs();
  setActAs(null);
  // Codeberg #4: re-resolve as the real admin now the grant is gone.
  clearIdentity();
  if (s) {
    try {
      await apiClient.delete(`/admin/impersonation/${s.id}`);
    } catch {
      // idempotent — grant may already be expired/gone
    }
  }
}
