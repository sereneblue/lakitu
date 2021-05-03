import type { AWSRegion, Response } from './types';

const login = async (password: string): Promise<Response> => {
	const formData = new URLSearchParams();
	formData.append('password', password);

	let res = await fetch(window.location.origin + '/login', {
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

const pingAWS = async (regions: string[]): Promise<Response> => {
	let res = await fetch(window.location.origin + '/ping-aws', {
		method: 'POST',
		credentials: 'include',
		headers: {
			Accept: 'application/json',
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({
			regions: regions
		})
	});

	if (res.ok) {
		return await res.json();
	}

	return {
		success: false,
		message: 'Unable to test latency'
	};
};

const verifyCredentials = async (accessKey: string, secretKey: string): Promise<Response> => {
	const formData = new URLSearchParams();
	formData.append('accessKey', accessKey);
	formData.append('secretKey', secretKey);

	let res = await fetch(window.location.origin + '/verify-creds', {
		method: 'POST',
		credentials: 'include',
		body: formData
	});

	if (res.ok) {
		return await res.json();
	}

	return {
		success: false,
		message: 'Credentials are incorrect'
	};
};

export { login, pingAWS, verifyCredentials };
