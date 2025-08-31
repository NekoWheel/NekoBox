import {toast, type ToastOptions} from 'vue3-toastify';

const defaultOptions: ToastOptions & { clearOnUrlChange: boolean } = {
    "theme": "auto",
    "position": "top-center",
    "pauseOnHover": false,
    "autoClose": 1500,
    "hideProgressBar": true,
    "transition": "slide",
    "clearOnUrlChange": false,
}

export function ToastInfo(message: string) {
    toast(message,
        {
            "type": "info",
            ...defaultOptions,
        }
    );
}

export function ToastSuccess(message: string) {
    toast(message,
        {
            "type": "success",
            ...defaultOptions,
        }
    );
}

export function ToastError(message: string) {
    toast(message,
        {
            "type": "error",
            ...defaultOptions,
        }
    );
}