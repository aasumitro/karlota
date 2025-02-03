import { AxiosError } from 'axios'
import { toast } from '@/hooks/use-toast'
import axios from "axios";
import {useAuthStore} from "@/states/auth-store";
import {isJWTExpired} from "@/lib/jwt";
import {Profile} from "@/types/user";

const host = "localhost:8000" // ${window.location.host} << use this in prod
const REST_URL = `${window.location.protocol}//${host}/api/v1`
const WS_URL = `ws://${host}/api/v1`

export const API_ENDPOINT = {
  ACCOUNT: {
    AUTH: {
      LOGIN: `/auth/login`,
      REGISTER: `/auth/register`,
      FORGOT_PASSWORD: `/auth/forgot-password`,
      RESET_PASSWORD: `/auth/reset-password`,
      REFRESH_TOKEN: `${REST_URL}/auth/refresh-token`,
    },
    PROFILE: `/account`,
    UPDATE_PASSWORD: `/account/password`,
  },
  MESSAGING: {
    CONTACT:  `/chats/contacts`,
    NEW_CONVERSATION:  `/chats`,
    WEBSOCKET: `${WS_URL}/chats`,
  }
}

export const HTTP_STATUS_CODE = {
  OK: 200,
  CREATED: 201,
  NOT_MODIFIED: 304,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  UNPROCESSABLE_ENTITY: 422,
  INTERNAL_SERVER_ERROR: 500
}

export interface HttpResponse<D , E> {
  data: D | null;
  errors: E | null; // errors can be string, array of objects, array of strings.
}

export interface ErrorResponse {
  errors: string | string[] | { [key: string]: string[] };
}

function isErrorResponse(error: unknown): error is ErrorResponse {
  return (
    typeof error === "object" &&
    error !== null &&
    "errors" in error &&
    (typeof (error as ErrorResponse).errors === "string" ||
      typeof (error as ErrorResponse).errors === "object")
  );
}

export function handleServerError(error: unknown) {
  let errMsg = "Something went wrong!";

  if (
    error &&
    typeof error === "object" &&
    "status" in error &&
    Number(error.status) === 204
  ) {
    errMsg = "Content not found.";
  }

  if (isErrorResponse(error)) {
    const errors = error.errors;
    if (typeof errors === "string") {
      errMsg = errors;
    } else if (typeof errors === "object") {
      // Flatten error messages into a single string
      errMsg = Object.entries(errors)
        .map(([key, messages]) => `${key}: ${messages.join(", ")}`)
        .join("; ");
    }
  }

  // Handle Axios errors
  if (error instanceof AxiosError) {
    const responseData = error.response;
    // If `title` exists in the response data, use it as the error message
    if (responseData?.data?.title) {
      errMsg = responseData.data.title;
    }
  }

  // Show the error message in a toast notification
  toast({ variant: "destructive", title: errMsg });
}

export const api = axios.create({
  baseURL: REST_URL,
  timeout: 3000,
  headers: {"Content-Type": "application/json"},
});

// Automatically add Authorization header to requests
api.interceptors.request.use(async  (config) => {
  let doRefresh = false;
  const auth = useAuthStore.getState().auth;

  if (auth.accessToken) {
    config.headers.Authorization = `Bearer ${auth.accessToken}`;
    if (isJWTExpired(auth.accessToken) && auth.refreshToken) {
      doRefresh = true;
    }
  }

  if (doRefresh && !isJWTExpired(auth.refreshToken)) {
    try {
      const response = await axios.post(API_ENDPOINT.ACCOUNT.AUTH.REFRESH_TOKEN,
        {}, {headers: {"X-REFRESH-TOKEN": auth.refreshToken}});
      if (response.status === HTTP_STATUS_CODE.CREATED) {
        const { access_token, user } = response.data.data;
        auth.setAccessToken(access_token.str);
        auth.setUser(user as Profile);
        config.headers.Authorization = `Bearer ${access_token.str}`;
      }
    } catch {
      auth.reset();
    }
  }

  return config;
});
