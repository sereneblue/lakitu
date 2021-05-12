<script context="module">
	export async function load({ page, fetch, session, context }) {
		const res = await fetch('http://' + page.host + '/setup/first-run');
		const data = await res.json();

		if (data.success) {
			return {
				status: 302,
				redirect: '/setup'
			};
		}

		return {
			status: 302,
			redirect: '/login'
		};
	}
</script>
