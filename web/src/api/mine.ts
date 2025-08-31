import axios from 'axios';

export interface MineQuestionItem {
    id: number;
    createdAt: Date;
    content: string;
    isAnswered: boolean;
    isPrivate: boolean;
}

export interface MineQuestions {
    total: string;
    cursor: string;
    questions: MineQuestionItem[];
}

export function mineQuestions(cursor: string, pageSize: number) {
    return axios.get<MineQuestions, MineQuestions>('/mine/questions', {
        params: {
            cursor,
            pageSize,
        }
    });
}

export interface AnswerQuestionRequest {
    answer: string;
    images: File[];
    recaptcha: string;
}

export function answerQuestion(questionID: number | string, data: AnswerQuestionRequest) {
    return axios.put<string, string>(`/mine/questions/${questionID}/answer`, data, {
        headers: {
            'Content-Type': 'multipart/form-data',
        }
    });
}

export function deleteQuestion(questionID: number | string) {
    return axios.delete<string, string>(`/mine/questions/${questionID}`);
}

export function setQuestionVisible(questionID: number | string, isVisible: boolean) {
    return axios.put<string, string>(`/mine/questions/${questionID}/visible`, {
        visible: isVisible,
    });
}

export interface MineProfile {
    email: string;
    name: string;
}

export function getMineProfile() {
    return axios.get<MineProfile, MineProfile>('/mine/settings/profile');
}

export interface UpdateProfileRequest {
    name: string;
    oldPassword?: string;
    newPassword?: string;
}

export function updateMineProfile(data: UpdateProfileRequest) {
    return axios.put<string, string>('/mine/settings/profile', data);
}