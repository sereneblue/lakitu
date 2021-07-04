import type { Response } from './types';

const login = async (password: string): Promise<Response> => {
	const formData = new URLSearchParams();
	formData.append('password', password);

	let res = await fetch(window.location.origin + '/session/login', {
		method: 'POST',
		credentials: 'include',
		body: formData
	});

	if (res.ok) {
		return await res.json();
	}

	return {
		success: false,
		message: 'Password is incorrect'
	};
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

	if (res.ok) {
		return await res.json();
	}
};

export { login, updatePassword };
