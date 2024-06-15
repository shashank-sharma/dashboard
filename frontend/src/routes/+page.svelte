<script lang="ts">
	import type { PageData } from "./$types";
	import { Button } from "attractions";
	export let data: PageData;
</script>

<h1 class="text-3xl font-bold underline text-red-400">Hello world!</h1>
<Button>Hello</Button>

{#if data.tasks}
	<nav class="list-nav card p-4 max-w-[400px] mx-auto space-y-4 my-4">
		<h2 class="h2 ml-4">TODOS</h2>
		<ul>
			{#each data.tasks.items as task, i}
				<li>
					<a href="/elements/lists">
						<span class="badge bg-secondary-500">🤟</span>
						<span class="flex-auto">{task.title}</span>
						<span class="">{task.description}</span>
					</a>
				</li>
			{/each}
		</ul>
	</nav>
{/if}
