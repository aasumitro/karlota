import Cookies from 'js-cookie'
import { create } from 'zustand'
import { Profile } from "@/types/user";

const ACCESS_TOKEN = 'KATS'
const REFRESH_TOKEN = 'KFTS'
const USER_PROFILE = 'KUPS'

interface AuthState {
  auth: {
    accessToken: string
    setAccessToken: (accessToken: string) => void

    refreshToken: string
    setRefreshToken: (refreshToken: string) => void

    user: Profile | null
    setUser: (user: Profile | null) => void

    reset: () => void
  }
}

export const useAuthStore = create<AuthState>()((set) => {
  const accessTokenCookieState = Cookies.get(ACCESS_TOKEN)
  const initAccessTokenToken = accessTokenCookieState ? JSON.parse(accessTokenCookieState) : ''

  const refreshTokenCookieState = Cookies.get(REFRESH_TOKEN)
  const initRefreshTokenToken = refreshTokenCookieState ? JSON.parse(refreshTokenCookieState) : ''

  const userProfileCookieState = Cookies.get(USER_PROFILE)
  const initUserProfile = userProfileCookieState ? JSON.parse(userProfileCookieState) : ''

  return {
    auth: {
      accessToken: initAccessTokenToken,
      setAccessToken: (accessToken) =>
        set((state) => {
          Cookies.set(ACCESS_TOKEN, JSON.stringify(accessToken), {
            sameSite: 'None',
            secure: true,
          })
          return { ...state, auth: { ...state.auth, accessToken } }
        }),

      refreshToken: initRefreshTokenToken,
      setRefreshToken: (refreshToken) =>
        set((state) => {
          Cookies.set(REFRESH_TOKEN, JSON.stringify(refreshToken), {
            sameSite: 'None',
            secure: true,
          })
          return { ...state, auth: { ...state.auth, refreshToken } }
        }),

      user: initUserProfile,
      setUser: (user) =>
        set((state) => {
          Cookies.set(USER_PROFILE, JSON.stringify(user), {
            sameSite: 'None',
            secure: true,
          })
          return { ...state, auth: { ...state.auth, user } }
        }),

      reset: () =>
        set((state) => {
          Cookies.remove(ACCESS_TOKEN)
          Cookies.remove(REFRESH_TOKEN)
          Cookies.remove(USER_PROFILE)
          return {
            ...state,
            auth: { ...state.auth, user: null, accessToken: '', refreshToken: '' },
          }
        }),
    },
  }
})
