import type {Response} from "./types";

const getJobStatus = async (jobId: number): Promise<Response> => {
	let res;

	if (jobId > 0) {
		res = await fetch(window.location.origin + '/jobs/' + jobId);
	} else {
		res = await fetch(window.location.origin + '/jobs');
	}

	if (res.ok) {
		return res.json();
	}

	return {
		success: false,
		message: "Unable to get job status"
	} 
}

export { getJobStatus };
