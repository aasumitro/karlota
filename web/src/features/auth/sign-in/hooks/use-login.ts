import {Tokens} from "@/types/user";
import {api, API_ENDPOINT, ErrorResponse, HttpResponse} from "@/lib/api";
import {useMutation} from "@tanstack/react-query";
import {AxiosError} from "axios";

export function useLogin() {
  const login = async (credentials: {
    email: string;
    password: string;
  }): Promise<HttpResponse<Tokens, ErrorResponse>> => {
    try {
      const response = await api.post<
        HttpResponse<Tokens, ErrorResponse
      >>(API_ENDPOINT.ACCOUNT.AUTH.LOGIN, credentials);
      return response.data;
    } catch (error: unknown) {
      if (error instanceof AxiosError && error.response) {
        const errorData = error.response.data;
        if (errorData && typeof errorData === "object") {
          throw errorData;
        }
      }
      throw new Error("Unexpected error occurred");
    }
  };

  return useMutation({mutationFn: login});
}
