import axios from "axios";

const instance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5
});

//axios.interceptors.request.use(function (config) {
//    const token = sessionStorage.getItem('token');
//	if (!token) return config;
//    config.headers.Authorization = "Bearer " + token;
//    return config;
//});

const updateToken = () => {
	instance.defaults.headers.common['Authorization'] = 'Bearer ' + sessionStorage.getItem('token');
}

export {
	instance as axios,
	updateToken as updateToken,
}
