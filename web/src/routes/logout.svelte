<script context="module">
	export async function load({ page, fetch, session, context }) {
		const res = await fetch('http://' + page.host + '/session/logout', {
			credentials: 'include'
		});

		await res.json();
		
		return {
			status: 302,
			redirect: '/login'
		};
	}
</script>
