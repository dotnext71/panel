import { describe, expect, it, vi, beforeEach } from "vitest";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";

const mockGet = vi.fn();
vi.mock("../../../apiClient", () => ({
  apiClient: { get: (url: string) => mockGet(url) },
}));

import { SupervisedWorkersCard } from "./SupervisedWorkersCard";

describe("SupervisedWorkersCard", () => {
  beforeEach(() => vi.clearAllMocks());

  it("renders supervised Go worker loops and running count", async () => {
    mockGet.mockResolvedValue({
      data: {
        total: 2,
        data: [
          { name: "reconciler", running: true, restarts: 0, last_start: "2026-11-18T00:00:00Z" },
          { name: "audit-consumer", running: true, restarts: 1, last_error: "" },
        ],
      },
    });
    const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(
      <QueryClientProvider client={qc}>
        <SupervisedWorkersCard />
      </QueryClientProvider>,
    );
    expect(await screen.findByText("reconciler")).toBeInTheDocument();
    expect(await screen.findByText("audit-consumer")).toBeInTheDocument();
    expect(await screen.findByText("2/2 running")).toBeInTheDocument();
  });
});
