import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "@fontsource/inter/400.css";
import "@fontsource/inter/500.css";
import "@fontsource/inter/600.css";
import "@fontsource/inter/700.css";
import "antd/dist/reset.css";
import "./global.css";

import "./i18n";
import { message, notification } from "antd";
import App from "./App";
import { registerServiceWorker } from "./lib/registerServiceWorker";
import { setPageTitle } from "./lib/pageTitle";

// tab title = "<hostname> | dotNext" so multiple instances are distinguishable
// at a glance. Set before render (independent of the React tree).
setPageTitle();

message.config({ top: 24, duration: 3, maxCount: 3 });
notification.config({ placement: "top", top: 24, duration: 4.5, maxCount: 3 });

window.addEventListener("vite:preloadError", () => {
  const last = Number(sessionStorage.getItem("dotnext-chunk-reload-at") || 0);
  if (Date.now() - last > 10_000) {
    sessionStorage.setItem("dotnext-chunk-reload-at", String(Date.now()));
    window.location.reload();
  }
});

const rootEl = document.getElementById("root");
if (!rootEl) {
  // This would mean index.html was modified incorrectly. Fail loud so
  // nobody ships a page that renders nothing.
  throw new Error("root element missing; check index.html");
}

createRoot(rootEl).render(
  <StrictMode>
    <App />
  </StrictMode>,
);

// #434: register the app service worker so the panel is an installable PWA.
registerServiceWorker();
