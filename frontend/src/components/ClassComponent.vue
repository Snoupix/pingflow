<script setup lang="ts">
import SpellsComponent from "@/components/SpellsComponent.vue";

import { useWebsocket } from "@/stores/websocket";
import type { Class } from "@/types/api";

const props = defineProps<{
	_class: Class;
	ws: ReturnType<typeof useWebsocket>;
}>();
</script>

<template>
	<h2>{{ props._class.name }}</h2>
	<section class="class_wrapper">
		<div class="grid">
			<template v-if="props._class.spellcasting != undefined">
				<span>Spells</span>
				<div>
					<SpellsComponent
						:ws="props.ws"
						:spells_endpoint="props._class.spells" />
				</div>

				<template v-for="({ name, desc }, i) in props._class.spellcasting!.info" :key="i">
					<span>{{ name }}</span>
					<span v-if="desc.length == 1">{{ desc[0] }}</span>
					<div v-else class="inner-grid">
						<p v-for="(line, i) in desc" :key="i">- {{ line }}</p>
					</div>
				</template>
			</template>
		</div>
	</section>
</template>

<style scoped lang="scss">
$border_radius: 0.5rem;

h2 {
	color: var(--color-text);
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
		border-radius: $border_radius;
		display: grid;
		grid-template-columns: 1fr 6fr;
		text-align: left;

		> * {
			padding: 15px;
			border: 1px dashed var(--color-border);
			display: flex;
			justify-content: center;
			align-items: center;

			button {
				border: 1px solid var(--color-border);
				background: var(--color-background-soft);
				padding: 2vh 5vw;
				color: var(--color-text);
				transition: all var(--transition-duration);
				border-radius: $border_radius;

				&:hover {
					border: 1px solid var(--color-background-soft);
					background: var(--color-border);
				}
			}
		}

		> :nth-child(1),
		> :nth-child(2) {
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
		> :nth-last-child(1),
		> :nth-last-child(2) {
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
