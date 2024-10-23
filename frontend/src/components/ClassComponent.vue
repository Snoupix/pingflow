<script setup lang="ts">
// import { ref } from 'vue';

import type { Class } from '@/types/api';

const props = defineProps<{
	_class: Class
}>();

// TODO: subtype
// const api_resp = ref<unknown | null>(null);
// const api_err = ref<unknown | null>(null);
</script>

<template>
	<h2>{{ props._class.name }}</h2>
	<section class="class_wrapper">
		<div class="grid">
			<template v-if="props._class.spellcasting != undefined">
				<span>Fetch spells</span>
				<span>{{ props._class.spells }}</span>
				<template v-for="{ name, desc }, i in props._class.spellcasting!.info" :key="i">
					<span>{{ name }}</span>
					<span v-if="desc.length == 1">{{ desc[0] }}</span>
					<div v-else class="inner-grid">
						<p v-for="line, i in desc" :key="i">- {{ line }}</p>
					</div>
				</template>
			</template>
		</div>
		<!--<button
			v-for="resp in api_resp.results"
			:key="resp.index"
			@click="FetchAPI({ endpoint: resp.url.replace(resp.index, ''), parameters: resp.index })">
			{{ resp.name }}
		</button>-->
	</section>
</template>

<style scoped lang="scss">
h2 {
	width: 100%;
	text-align: center;
}

.class_wrapper {
	display: flex;
	flex-direction: row;
	justify-content: center;
	align-items: center;
	gap: 10px;
	width: 100%;
	text-align: center;
	padding: 5vh 5vw;

	.grid {
		border-radius: .5rem;
		display: grid;
		grid-template-columns: 1fr 6fr;
		text-align: left;

		> * {
			padding: 15px;
			border: 1px dashed var(--color-border);
			display: flex;
			justify-content: center;
			align-items: center;
		}

		> :nth-child(1), > :nth-child(2) {
			border-top: none;
		}
		> :nth-child(odd) {
			border-left: none;
			text-align: center;
		}
		> :nth-child(even) {
			border-right: none;
			justify-content: left;
		}
		> :nth-last-child(1), > :nth-last-child(2) {
			border-bottom: none;
		}

		.inner-grid {
			display: flex;
			flex-direction: column;

			> * {
				width: 100%;
			}
		}
	}

	& > button {
		border: 1px solid var(--color-border);
		background: transparent;
		display: flex;
		justify-content: center;
		align-items: center;
		width: 100%;
		height: 10vh;
		cursor: pointer;
		color: var(--color-text);
		transition: all 0.3s;

		&:hover {
			border-color: var(--color-border-hover);
			background: var(--color-background-soft);
			color: white;
		}
	}
}
</style>
