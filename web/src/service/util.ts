import { Notification } from '../components/base/';
import type { NotificationType } from '../service/types';

export function notify(
	target: HTMLElement,
	notifyType: NotificationType,
	message: string,
	duration: number = 2500
) {
	new Notification({
		target,
		hydrate: true,
		props: {
			message,
			type: notifyType,
			duration
		}
	});
}
