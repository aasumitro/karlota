import {api, API_ENDPOINT, ErrorResponse, HttpResponse} from "@/lib/api";
import {Profile} from "@/types/user";
import {useMutation, useSuspenseQuery} from "@tanstack/react-query";
import {AxiosError} from "axios";

export function useProfile() {
  const profile = async (): Promise<HttpResponse<Profile, ErrorResponse>> => {
    try {
      const response = await api.get<
        HttpResponse<Profile, ErrorResponse
        >>(API_ENDPOINT.ACCOUNT.PROFILE);
      return response.data;
    } catch (error: unknown) {
      if (error instanceof AxiosError && error.response) {
        const errorData = error.response.data;
        if (errorData && typeof errorData === "object") {
          throw errorData;
        }
        throw error
      }
      throw new Error("Unexpected error occurred");
    }
  };

  return useSuspenseQuery({ queryKey: ['profile'], queryFn: profile })
}

export function useUpdateDisplayName() {
  const profile = async (
    data: {
      display_name: string
    }
  ): Promise<HttpResponse<Profile, ErrorResponse>> => {
    try {
      const response = await api.patch<
        HttpResponse<Profile, ErrorResponse
        >>(API_ENDPOINT.ACCOUNT.PROFILE, data);
      return response.data;
    } catch (error: unknown) {
      if (error instanceof AxiosError && error.response) {
        const errorData = error.response.data;
        if (errorData && typeof errorData === "object") {
          throw errorData;
        }
        throw error
      }
      throw new Error("Unexpected error occurred");
    }
  };

  return useMutation({ mutationFn: profile })
}

export function useUpdatePassword() {
  const profile = async (
    data: {
      old_password: string
      new_password: string
    }
  ): Promise<HttpResponse<Profile, ErrorResponse>> => {
    try {
      const response = await api.patch<
        HttpResponse<Profile, ErrorResponse
        >>(API_ENDPOINT.ACCOUNT.UPDATE_PASSWORD, data);
      return response.data;
    } catch (error: unknown) {
      if (error instanceof AxiosError && error.response) {
        const errorData = error.response.data;
        if (errorData && typeof errorData === "object") {
          throw errorData;
        }
        throw error
      }
      throw new Error("Unexpected error occurred");
    }
  };

  return useMutation({ mutationFn: profile })
}