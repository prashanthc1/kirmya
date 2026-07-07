import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import {
  ApiError,
  api,
  fetchCsrfToken,
  getAccessToken,
  request,
  setAccessToken,
} from "@/lib/api/client";

/** Builds a Response-like object for the stubbed fetch. */
function jsonResponse(status: number, body: unknown): Response {
  return {
    ok: status >= 200 && status < 300,
    status,
    statusText: `status ${status}`,
    json: async () => body,
    text: async () => (body === undefined ? "" : JSON.stringify(body)),
  } as unknown as Response;
}

describe("api client", () => {
  beforeEach(() => {
    setAccessToken(null);
    // Clear cookies between tests.
    document.cookie =
      "csrf_token=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/";
    vi.restoreAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("unwraps the { data } envelope on success", async () => {
    const fetchMock = vi
      .spyOn(globalThis, "fetch")
      .mockResolvedValue(
        jsonResponse(200, { data: { id: "u1", name: "Ada" } }),
      );

    const result = await api.get<{ id: string; name: string }>("/users/me");

    expect(result).toEqual({ id: "u1", name: "Ada" });
    expect(fetchMock).toHaveBeenCalledTimes(1);
    const url = fetchMock.mock.calls[0][0];
    expect(url).toBe("/api/v1/users/me");
  });

  it("attaches the Bearer token when set", async () => {
    setAccessToken("tok-123");
    const fetchMock = vi
      .spyOn(globalThis, "fetch")
      .mockResolvedValue(jsonResponse(200, { data: true }));

    await api.get("/users/me");

    const init = fetchMock.mock.calls[0][1] as RequestInit;
    const headers = new Headers(init.headers);
    expect(headers.get("Authorization")).toBe("Bearer tok-123");
    expect(init.credentials).toBe("include");
  });

  it("sends the X-CSRF-Token header from the cookie on mutating requests", async () => {
    document.cookie = "csrf_token=csrf-abc; path=/";
    const fetchMock = vi
      .spyOn(globalThis, "fetch")
      .mockResolvedValue(jsonResponse(200, { data: { ok: true } }));

    await api.post("/auth/logout");

    const init = fetchMock.mock.calls[0][1] as RequestInit;
    const headers = new Headers(init.headers);
    expect(headers.get("X-CSRF-Token")).toBe("csrf-abc");
  });

  it("does not send a CSRF header on safe GET requests", async () => {
    document.cookie = "csrf_token=csrf-abc; path=/";
    const fetchMock = vi
      .spyOn(globalThis, "fetch")
      .mockResolvedValue(jsonResponse(200, { data: 1 }));

    await api.get("/jobs");

    const init = fetchMock.mock.calls[0][1] as RequestInit;
    const headers = new Headers(init.headers);
    expect(headers.get("X-CSRF-Token")).toBeNull();
  });

  it("throws ApiError carrying the error envelope", async () => {
    vi.spyOn(globalThis, "fetch").mockResolvedValue(
      jsonResponse(422, {
        error: { code: "validation_error", message: "bad email" },
      }),
    );

    await expect(api.post("/auth/register", {})).rejects.toMatchObject({
      name: "ApiError",
      status: 422,
      code: "validation_error",
      message: "bad email",
    });
  });

  it("refreshes on 401 then retries the original request once", async () => {
    const fetchMock = vi
      .spyOn(globalThis, "fetch")
      // 1st: original request -> 401
      .mockResolvedValueOnce(
        jsonResponse(401, {
          error: { code: "unauthorized", message: "expired" },
        }),
      )
      // 2nd: refresh -> new token
      .mockResolvedValueOnce(
        jsonResponse(200, { data: { access_token: "fresh-tok" } }),
      )
      // 3rd: retried original -> success
      .mockResolvedValueOnce(jsonResponse(200, { data: { id: "u1" } }));

    const result = await api.get<{ id: string }>("/users/me");

    expect(result).toEqual({ id: "u1" });
    expect(fetchMock).toHaveBeenCalledTimes(3);
    expect(getAccessToken()).toBe("fresh-tok");
    // The refresh call hit the refresh endpoint.
    expect(fetchMock.mock.calls[1][0]).toBe("/api/v1/auth/refresh");
  });

  it("clears the token and surfaces 401 when refresh fails", async () => {
    setAccessToken("stale");
    vi.spyOn(globalThis, "fetch")
      .mockResolvedValueOnce(
        jsonResponse(401, {
          error: { code: "unauthorized", message: "expired" },
        }),
      )
      .mockResolvedValueOnce(
        jsonResponse(401, {
          error: { code: "unauthorized", message: "no cookie" },
        }),
      );

    await expect(api.get("/users/me")).rejects.toBeInstanceOf(ApiError);
    expect(getAccessToken()).toBeNull();
  });

  it("fetchCsrfToken returns the token from the envelope", async () => {
    vi.spyOn(globalThis, "fetch").mockResolvedValue(
      jsonResponse(200, { data: { csrf_token: "xyz" } }),
    );

    await expect(fetchCsrfToken()).resolves.toBe("xyz");
  });

  it("returns undefined for 204 responses", async () => {
    vi.spyOn(globalThis, "fetch").mockResolvedValue(
      jsonResponse(204, undefined),
    );
    await expect(
      request("/auth/logout", { method: "POST" }),
    ).resolves.toBeUndefined();
  });
});
