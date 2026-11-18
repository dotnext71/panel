// muiTheme.ts — dotNext AntD ConfigProvider design system.
//
// Dark-first SaaS palette: obsidian slate canvas (#0B0F17), layered card
// surfaces (#111827 / #182235), crisp 1px borders (#1E293B), and an electric
// cyan accent (#06B6D4 / #0891B2). Light mode pairs #F8FAFC canvas with pure
// white cards and slate hairline borders. Operator Look-and-Feel overrides from
// /branding merge cleanly atop these defaults.
import { useMemo } from "react";
import { theme } from "antd";
import type { ConfigProviderProps } from "antd";

import type { ThemeMode } from "./theme/ThemeModeContext";

type LookAndFeel = {
  fontSizePx?: number;
  accentColor?: string;
  seedColors?: Record<string, string>;
  chrome?: { bg?: string; container?: string; text?: string };
};

const useMuiTheme = (
  mode: ThemeMode,
  opts?: LookAndFeel,
): ConfigProviderProps => {
  const fontSizePx = opts?.fontSizePx ?? 14;
  const seed = opts?.seedColors ?? {};
  const chrome = opts?.chrome ?? {};
  const isDark = mode === "dark";

  const accent =
    opts?.accentColor || (isDark ? "#06B6D4" : "#0891B2");
  const primary =
    seed.colorPrimary || (isDark ? "#06B6D4" : "#0891B2");

  return useMemo<ConfigProviderProps>(
    () => ({
      theme: {
        algorithm: isDark ? theme.darkAlgorithm : theme.defaultAlgorithm,
        token: {
          fontSize: fontSizePx,
          borderRadius: 8,
          borderRadiusLG: 12,
          borderRadiusSM: 6,
          colorPrimary: primary,
          colorSuccess: seed.colorSuccess || (isDark ? "#10B981" : "#059669"),
          colorWarning: seed.colorWarning || (isDark ? "#F59E0B" : "#D97706"),
          colorError: seed.colorError || (isDark ? "#F43F5E" : "#E11D48"),
          colorInfo: seed.colorInfo || (isDark ? "#38BDF8" : "#0284C7"),
          colorLink: seed.colorLink || primary,
          colorBgLayout: chrome.bg || (isDark ? "#0B0F17" : "#F8FAFC"),
          colorBgContainer:
            chrome.container || (isDark ? "#111827" : "#FFFFFF"),
          colorBgElevated: isDark ? "#182235" : "#FFFFFF",
          colorBorder: isDark ? "#1E293B" : "#E2E8F0",
          colorBorderSecondary: isDark ? "#1E293B" : "#E2E8F0",
          ...(chrome.text
            ? { colorTextBase: chrome.text }
            : isDark
              ? { colorText: "#F1F5F9", colorTextSecondary: "#94A3B8" }
              : { colorText: "#0F172A", colorTextSecondary: "#475569" }),
          fontFamily:
            "Inter, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif",
          fontFamilyCode:
            "'JetBrains Mono', ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace",
        },
        components: {
          Layout: {
            triggerBg: "transparent",
            bodyBg: chrome.bg || (isDark ? "#0B0F17" : "#F8FAFC"),
            siderBg: chrome.bg || (isDark ? "#0E1420" : "#FFFFFF"),
            headerBg: "transparent",
          },
          Card: {
            headerBg: "transparent",
            borderRadiusLG: 12,
          },
          Table: {
            headerBg: isDark ? "#141D2E" : "#F8FAFC",
            headerColor: isDark ? "#94A3B8" : "#475569",
            rowHoverBg: isDark ? "rgba(6, 182, 212, 0.06)" : "#F1F5F9",
            borderColor: isDark ? "#1E293B" : "#E2E8F0",
          },
          Menu: isDark
            ? {
                itemBg: "transparent",
                darkItemBg: "transparent",
                darkItemColor: "#CBD5E1",
                darkItemSelectedBg: "rgba(6, 182, 212, 0.14)",
                darkItemSelectedColor: accent,
                darkItemHoverBg: "rgba(148, 163, 184, 0.09)",
                darkItemHoverColor: "#F8FAFC",
                itemBorderRadius: 8,
              }
            : {
                itemBg: "transparent",
                itemColor: "#334155",
                itemSelectedBg: "rgba(8, 145, 178, 0.10)",
                itemSelectedColor: accent,
                itemHoverBg: "#F1F5F9",
                itemHoverColor: "#0F172A",
                itemBorderRadius: 8,
              },
          Tabs: {
            itemSelectedColor: accent,
            inkBarColor: accent,
            itemHoverColor: accent,
            itemActiveColor: accent,
          },
          Button: {
            controlHeight: 36,
            borderRadius: 8,
          },
          Input: {
            controlHeight: 36,
            borderRadius: 8,
          },
          Select: {
            controlHeight: 36,
            borderRadius: 8,
          },
        },
      },
    }),
    [isDark, mode, fontSizePx, accent, primary, JSON.stringify(seed), JSON.stringify(chrome)],
  );
};

export default useMuiTheme;
