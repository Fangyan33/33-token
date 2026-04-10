export type ApiClientConfig = {
  baseUrl?: string;
};

export class ApiClient {
  private readonly baseUrl: string;

  constructor(config: ApiClientConfig = {}) {
    this.baseUrl = config.baseUrl ?? "/api/v1";
  }

  async get(path: string): Promise<Response> {
    return fetch(`${this.baseUrl}${path}`, {
      credentials: "include"
    });
  }
}
