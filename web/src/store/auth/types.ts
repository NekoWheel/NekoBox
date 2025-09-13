import type {UserProfile} from "@/api/auth.ts";

export interface AuthState {
    isSignedIn: boolean;

    profile: UserProfile;
    sessionID: string;
}