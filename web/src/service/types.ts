export type AWSRegion = {
	id: string;
	name: string;
};

export type Response = {
	success: boolean;
	message?: string;
	data?: object;
};
