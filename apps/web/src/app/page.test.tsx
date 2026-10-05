import { renderToStaticMarkup } from "react-dom/server";
import { describe, expect, it } from "vitest";
import Home from "./page";

describe("Home", () => {
  it("identifies the local testing product", () => {
    const output = renderToStaticMarkup(<Home />);

    expect(output).toContain("<h1>Device Force</h1>");
    expect(output).toContain("local Agent is ready");
  });
});
