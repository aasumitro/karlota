import {api, API_ENDPOINT, ErrorResponse, HttpResponse} from "@/lib/api";
import {Profile} from "@/types/user";
import {AxiosError} from "axios";
import {useQuery} from "@tanstack/react-query";

export function useContact() {
  const contacts = async (): Promise<HttpResponse<Profile[], ErrorResponse>> => {
    try {
      const response = await api.get<
        HttpResponse<Profile[], ErrorResponse
        >>(API_ENDPOINT.MESSAGING.CONTACT);
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

  return useQuery({ queryKey: ['contacts'], queryFn: contacts })
}