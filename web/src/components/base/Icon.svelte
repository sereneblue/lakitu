<script lang="typescript">
	import feather from 'feather-icons';

	export let icon: string = '';
	export let stroke: number = 2;
</script>

{@html feather.icons[icon].toSvg({ width: '', height: '', 'stroke-width': stroke })}
