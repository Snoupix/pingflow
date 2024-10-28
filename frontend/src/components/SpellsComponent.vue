<script setup lang="ts">
import { ref, watch } from "vue";

import SpellInfoComponent from "@/components/SpellInfo.vue";

import { useFetcher } from "@/stores/fetcher";
import { useWebsocket, type IWorkerConfig } from "@/stores/websocket";
import type { SpellInfo, SpellResp } from "@/types/api";

const props = defineProps<{
	spells_endpoint: string;
	ws: ReturnType<typeof useWebsocket>;
}>();

const { Fetch } = useFetcher();
const spells = ref<SpellResp | null>(null);
const spell_info = ref<SpellInfo | null>(null);
const processed_spells = ref<Map<number, Array<{ url: string; name: string }>>>(new Map());

watch(spells, s => {
	if (s == null) return;

	s.results.forEach(spell => {
		const array = processed_spells.value.get(spell.level);
		if (array == undefined) {
			processed_spells.value.set(spell.level, [{ ...spell }]);
			return;
		}

		array.push({ ...spell });
		processed_spells.value.set(spell.level, array);
	});
});

function FetchSpells(cfg: IWorkerConfig) {
	Fetch(props.ws, spells, cfg);
}

function FetchSpellInfo(cfg: IWorkerConfig) {
	Fetch(props.ws, spell_info, cfg);
}
</script>

<template>
	<div v-if="spell_info != null">
		<SpellInfoComponent :spell_info />
		<button class="close_info" @click="() => spell_info = null">Close</button>
	</div>
	<button
		v-else-if="spells == null"
		@click="FetchSpells({ endpoint: props.spells_endpoint, parameters: '' })">
		Fetch spells
	</button>
	<div v-else class="spell_grid">
		<div v-for="([level, names], i) in processed_spells" :key="i">
			<span>Level {{ level }}</span>
			<button
				class="spell"
				v-for="({ url, name }, j) in names"
				:key="j"
				@click="FetchSpellInfo({ endpoint: url, parameters: '' })">
				{{ name }}
			</button>
		</div>
	</div>
</template>

<style scoped lang="scss">
.close_info {
	background: var(--color-background);
	border: 1px solid var(--color-border);
	border-radius: 0.5rem;
	color: inherit;
	font: inherit;
	width: 100px;
	height: 35px;
	margin-top: 10px;
	transition: all var(--transition-duration);

	&:hover {
		background: var(--color-background-soft);
		border-color: var(--color-background);
	}
}

.spell_grid {
	display: flex;
	flex-direction: row;
	justify-content: center;
	align-items: start;
	text-align: center;
	max-height: 500px;
	overflow-y: scroll;

	> * {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		gap: 5px;
		border-right: 2px dashed var(--color-border);

		> * {
			padding: 10px 2vw;
		}

		& > :first-child {
			width: 100%;
			border-bottom: 5px dotted var(--color-border);
		}
	}

	& > :last-child {
		border-right: none;
	}

	.spell {
		background: none;
		border: none;
		color: inherit;
		font: inherit;
		border-color: var(--color-border);
		transition: all var(--transition-duration);

		&:hover {
			border-radius: 10rem;
			background: var(--color-background-soft);
		}
	}
}
</style>
