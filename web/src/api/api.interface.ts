import {AxiosRequestConfig, AxiosResponse} from "axios";

export const allowedApiMethods = ["get", "post", "put", "patch", "delete"];

export interface ApiClientMethod {
  (url: string, config?: AxiosRequestConfig): Promise<AxiosResponse<any>>;
}

export interface ApiClientMethodWithData {
  (url: string, data?: any, config?: AxiosRequestConfig): Promise<AxiosResponse<any>>;
}

export interface ApiClient {
  get: ApiClientMethod;
  post: ApiClientMethodWithData;
  patch: ApiClientMethodWithData;
  put: ApiClientMethodWithData;
  delete: ApiClientMethodWithData;
}
