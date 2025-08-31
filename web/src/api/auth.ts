import axios from 'axios'


export interface SignUpRequest {
    email: string;
    domain: string;
    name: string;
    password: string;
    repeatPassword: string;
    recaptcha: string;
}

export function signUp(data: SignUpRequest) {
    return axios.post<string, string>('/auth/sign-up', data)
}

export interface SignInRequest {
    email: string;
    password: string;
    recaptcha: string;
}

export interface SignInResponse {
    profile: UserProfile;
    sessionID: string;
}

export interface UserProfile {
    uid: string;
    name: string;
    domain: string;
}

export function signIn(data: SignInRequest) {
    return axios.post<SignInResponse, SignInResponse>('/auth/sign-in', data)
}
