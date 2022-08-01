import axios from 'axios'
import stub from "@/init/index";

enum Method {
    GET = "GET",
    POST = "POST",
    PUT = "PUT",
    DELETE = "DELETE"
}

const service = axios.create({
    baseURL: '/api',
    timeout: 5000,
});
service.interceptors.request.use(
    request => request,
    error => Promise.reject(error)
);

service.interceptors.response.use(
    response => {
        const error = response.data?.error
        if (error) {
            // stub.ref.antd.Modal.error({content: error.message, centered: true})
        }
        return response
    },
    error => Promise.reject(error)
);

function apiAxios(method: Method, url: string, params: any) {
    return new Promise((resolve, reject) => {
        service({
            method: method,
            url: url,
            data: method === 'POST' || method === 'PUT' ? params : null,
            params: method === 'GET' || method === 'DELETE' ? params : null
        }).then(res => {
            return resolve(res)
        }, error => {
            return reject(error)
        }).catch(error => reject(error))
    })
}

export default {
    get: (url: string, params?: any) => {
        return apiAxios(Method.GET, url, params)
    },
    post: (url: string, params?: any) => {
        return apiAxios(Method.POST, url, params)
    },
    put: (url: string, params?: any) => {
        return apiAxios(Method.PUT, url, params)
    },
    delete: (url: string, params?: any) => {
        return apiAxios(Method.DELETE, url, params)
    }
};