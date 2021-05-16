import type { AWSRegion, Response } from './types';

const completeSetup = async (accessKey, secretKey, region, password ): Promise<Response> => {
	const formData = new URLSearchParams();
	formData.append('accessKey', accessKey);
	formData.append('secretKey', secretKey);
	formData.append('password', password);
	formData.append('region', region);

	let res = await fetch(window.location.origin + '/setup/complete', {
		method: 'POST',
		credentials: 'include',
		headers: {
			Accept: 'application/json',
		},
		body: formData
	});

	if (res.ok) {
		return await res.json();
	}

	return {
		success: false,
		message: 'Unable to complete setup'
	};
};

const pingAWS = async (regions: string[]): Promise<Response> => {
	let res = await fetch(window.location.origin + '/aws/ping-test', {
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

	let res = await fetch(window.location.origin + '/aws/verify-creds', {
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

export { completeSetup, pingAWS, verifyCredentials };
