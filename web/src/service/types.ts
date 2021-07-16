export type AWSRegion = {
	id: string;
	name: string;
};

export type Response = {
	success: boolean;
	message?: string;
	data?: object;
};

export type NotificationType = 'warning' | 'info' | 'danger' | 'success';
