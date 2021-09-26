export type AWSRegion = {
	id: string;
	name: string;
};

export type Log = {
	id: number;
	event: string;
	errorDetails: string;
	timestamp: number;
	status: Status;
};

export type MachineData = {
	uuid: string;
	instanceType: string;
	name: string;
	region: string;
	size: number;
	status: MachineStatus;
}

export type MachineEvent = {
	event: MachineEventType;
	machine: MachineData;
}

export type ModalAction = {
	text: string;
	isDisabled?: boolean;
	hideCancel: boolean;
	func: Function;
};

export type Response = {
	success: boolean;
	message?: string;
	data?: object;
};

export const enum Status {
	PENDING = 0,
	COMPLETE = 1,
	ERROR = 2,
	CANCELED = 3,
};

export type MachineEventType = 'start' | 'stop' | 'increaseStorage' | 'changeRegion' | 'delete';
export type MachineStatus = 'offline' | 'online' | 'unavailable';
export type ModalType = 'warning' | 'info' | 'danger' | 'success' | 'default';
export type NotificationType = 'warning' | 'info' | 'danger' | 'success';

