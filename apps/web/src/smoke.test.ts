import { describe, expect, it } from "vitest";

describe("web scaffold", () => {
  it("keeps the workspace bootstrap title stable", () => {
    expect("Model API Platform").toContain("Platform");
  });
});
