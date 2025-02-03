import {api, API_ENDPOINT, ErrorResponse, HttpResponse} from "@/lib/api";
import {AxiosError} from "axios";
import {useMutation} from "@tanstack/react-query";

export function useConversation() {
  const newConversation = async (
    data: {
      name?: string
      message?: string
      members: {id?: number, display_name?: string}[]
    }
  ): Promise<HttpResponse<string, ErrorResponse>> => {
    try {
      const response = await api.post<
        HttpResponse<string, ErrorResponse
        >>(API_ENDPOINT.MESSAGING.NEW_CONVERSATION, data);
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

  return useMutation({ mutationFn: newConversation })
}