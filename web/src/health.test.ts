import { describe, expect, it } from "vitest";
import { isHealthy } from "./health";

describe("isHealthy", () => {
  it("accepts ok", () => expect(isHealthy({ status: "ok" })).toBe(true));
  it("rejects anything else", () => expect(isHealthy({ status: "down" })).toBe(false));
});
