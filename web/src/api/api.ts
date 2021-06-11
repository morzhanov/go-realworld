import axios, {AxiosResponse, AxiosRequestConfig} from "axios";

import {ApiClient} from "./api.interface";

export function createApiClient(baseURL = ""): ApiClient {
  return {
    get: (url, config: AxiosRequestConfig = {}): Promise<AxiosResponse> =>
      axios.get(url, {...config, baseURL}),
    post: (url, data, config: AxiosRequestConfig = {}): Promise<AxiosResponse> =>
      axios.post(url, data, {...config, baseURL}),
    patch: (url, data, config: AxiosRequestConfig = {}): Promise<AxiosResponse> =>
      axios.patch(url, data, {...config, baseURL}),
    put: (url, data, config: AxiosRequestConfig = {}): Promise<AxiosResponse> =>
      axios.put(url, data, {...config, baseURL}),
    delete: (url, data, config: AxiosRequestConfig = {}): Promise<AxiosResponse> =>
      axios.delete(url, {...config, baseURL, data}),
  };
}

export const api = createApiClient(process.env.API_URI);
