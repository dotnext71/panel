// Support links. SupportPage hides any CTA whose value is empty, so
// removing one here gracefully degrades to a "Coming soon" tag rather
// than a dead link.
export const SUPPORT_LINKS = {
  documentation: "https://github.com/dotnext71/panel#readme",
  githubIssues: "https://github.com/dotnext71/panel/issues",
  paidSupport: "",
  emergency: "",
} as const;

export const DIAGNOSTIC_EMAIL_RECIPIENT = "";
