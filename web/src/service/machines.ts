import type { Response } from './types';

const getRegions = async (): Promise<Response> => {
	let res = await fetch(window.location.origin + '/aws/regions', {
		credentials: 'include'
	});

	if (res.ok) {
		return await res.json();
	}

	return {
		success: false
	};
};

const getMachines = async (): Promise<Response> => {
	let res = await fetch(window.location.origin + '/machine/list', {
		credentials: 'include'
	});

	if (res.ok) {
		return await res.json();
	}

	return {
		success: false
	};
};

const getInstances = async (region: string): Promise<Response> => {
	const formData = new URLSearchParams();
	formData.append('region', region);

	let res = await fetch(window.location.origin + '/aws/gpu-instances', {
		method: 'POST',
		credentials: 'include',
		body: formData
	});

	if (res.ok) {
		return await res.json();
	}

	return {
		success: false
	};	
};

const getPricing = async (region: string): Promise<Response> => {
	const formData = new URLSearchParams();
	formData.append('region', region);

	let res = await fetch(window.location.origin + '/aws/pricing', {
		method: 'POST',
		credentials: 'include',
		body: formData
	});

	if (res.ok) {
		return await res.json();
	}

	return {
		success: false
	};
};

const createMachine = async (form: object): Promise<Response> => {
	const formData = new URLSearchParams();
	formData.append('region', form.region);
	formData.append('size', form.size.toString());
	formData.append('instanceType', form.instanceType);
	formData.append('streamOption', form.streamOption);
	formData.append('name', form.name);

	let res = await fetch(window.location.origin + '/machine/create', {
		method: 'POST',
		credentials: 'include',
		body: formData
	});

	if (res.ok) {
		return await res.json();
	}

	return {
		success: false
	};
}

const deleteMachine = async (uuid: string): Promise<Response> => {
	const formData = new URLSearchParams();
	formData.append('uuid', uuid);

	let res = await fetch(window.location.origin + '/machine/delete', {
		method: 'POST',
		credentials: 'include',
		body: formData
	});

	if (res.ok) {
		return await res.json();
	}

	return {
		success: false
	};
}

const resizeMachine = async (uuid: string, size: number): Promise<Response> => {
	const formData = new URLSearchParams();
	formData.append('uuid', uuid);
	formData.append('size', size.toString());

	let res = await fetch(window.location.origin + '/machine/resize', {
		method: 'POST',
		credentials: 'include',
		body: formData
	});

	if (res.ok) {
		return await res.json();
	}

	return {
		success: false
	};
}

const startMachine = async (uuid: string): Promise<Response> => {
	const formData = new URLSearchParams();
	formData.append('uuid', uuid);

	let res = await fetch(window.location.origin + '/machine/start', {
		method: 'POST',
		credentials: 'include',
		body: formData
	});

	if (res.ok) {
		return await res.json();
	}

	return {
		success: false
	};
}

const stopMachine = async (uuid: string): Promise<Response> => {
	const formData = new URLSearchParams();
	formData.append('uuid', uuid);

	let res = await fetch(window.location.origin + '/machine/stop', {
		method: 'POST',
		credentials: 'include',
		body: formData
	});

	if (res.ok) {
		return await res.json();
	}

	return {
		success: false
	};
}

const transferMachine = async (uuid: string, region: string): Promise<Response> => {
	const formData = new URLSearchParams();
	formData.append('uuid', uuid);
	formData.append('region', region);

	let res = await fetch(window.location.origin + '/machine/transfer', {
		method: 'POST',
		credentials: 'include',
		body: formData
	});

	if (res.ok) {
		return await res.json();
	}

	return {
		success: false
	};
}

export { 
	getInstances, 
	getMachines, 
	getPricing, 
	getRegions, 
	createMachine, 
	deleteMachine, 
	resizeMachine,
	startMachine,
	stopMachine, 
	transferMachine 
};
