import type { AxiosResponse } from "axios";

export class ApiError extends Error {
  constructor(
    message: string,
    public details: {
      response: AxiosResponse;
      expected: string;
      received: unknown;
      code?: string;
      details?: unknown;
    }
  ) {
    super(message);
    this.name = 'ApiError';
  }
}