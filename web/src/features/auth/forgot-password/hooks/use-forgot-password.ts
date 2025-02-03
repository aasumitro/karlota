import {api, API_ENDPOINT, ErrorResponse, HttpResponse} from "@/lib/api";
import {useMutation} from "@tanstack/react-query";
import {AxiosError} from "axios";

export function useRequestResetPasswordLink() {
  const resetPasswordLink = async (credentials: {
    email: string;
  }): Promise<HttpResponse<string, ErrorResponse>> => {
    try {
      const response = await api.post<
        HttpResponse<string, ErrorResponse
        >>(API_ENDPOINT.ACCOUNT.AUTH.FORGOT_PASSWORD, credentials);
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

  return useMutation({mutationFn: resetPasswordLink});
}
