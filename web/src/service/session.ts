import type { Response } from './types';

const getUserData = async (host: string): Promise<Response> => {
	let res = await fetch('http://' + host + '/session/user');

	if (res.ok) {
		return await res.json();
	}

	return {
		success: false,
		message: 'Unable to get user data'
	};
};

const login = async (password: string): Promise<Response> => {
	const formData = new URLSearchParams();
	formData.append('password', password);

	let res = await fetch(window.location.origin + '/session/login', {
		method: 'POST',
		credentials: 'include',
		body: formData
	});

	return await res.json();
};

const updatePassword = async (
	oldPwd: string,
	newPwd: string,
	confirmNewPwd: string
): Promise<Response> => {
	const formData = new URLSearchParams();
	formData.append('oldPwd', oldPwd);
	formData.append('newPwd', newPwd);
	formData.append('confirmNewPwd', confirmNewPwd);

	let res = await fetch(window.location.origin + '/session/change-password', {
		method: 'POST',
		credentials: 'include',
		body: formData
	});

	return await res.json();
};

const updatePreferences = async (
	accessKey: string,
	secretKey: string,
	defaultRegion: string
): Promise<Response> => {
	const formData = new URLSearchParams();
	formData.append('accessKey', accessKey);
	formData.append('secretKey', secretKey);
	formData.append('defaultRegion', defaultRegion);

	let res = await fetch(window.location.origin + '/session/update-preferences', {
		method: 'POST',
		credentials: 'include',
		body: formData
	});

	return await res.json();
};

export { getUserData, login, updatePassword, updatePreferences };
