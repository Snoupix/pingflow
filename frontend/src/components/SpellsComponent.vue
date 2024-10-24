<script setup lang="ts">
import { ref } from "vue";

import SpellInfoComponent from "@/components/SpellInfo.vue";

import type { IWorkerConfig } from "@/stores/websocket";
import type { API_OUT_T, Spell, SpellInfo } from "@/types/api";

const props = defineProps<{
	spells_endpoint: string;
	spells: Array<Spell> | null;
	spell_info: SpellInfo | null;
	FetchAPI: (config: IWorkerConfig & { _type: API_OUT_T }) => void;
	ClearSpellInfo: () => void;
}>();

const processed_spells = ref<Map<number, Array<{ url: string; name: string }>>>(new Map());

if (props.spells != null) {
	props.spells.forEach(spell => {
		const array = processed_spells.value.get(spell.level);
		if (array == undefined) {
			processed_spells.value.set(spell.level, [{ ...spell }]);
			return;
		}

		array.push({ ...spell });
		processed_spells.value.set(spell.level, array);
	});
}
</script>

<template>
	<div v-if="props.spell_info != null">
		<SpellInfoComponent :spell_info="props.spell_info" />
		<button class="close_info" @click="props.ClearSpellInfo()">Close</button>
	</div>
	<button
		v-else-if="props.spells == null"
		@click="props.FetchAPI({ endpoint: props.spells_endpoint, parameters: '', _type: 'spells' })">
		Fetch spells
	</button>
	<div v-else class="spell_grid">
		<div v-for="([level, names], i) in processed_spells" :key="i">
			<span>Level {{ level }}</span>
			<button
				class="spell"
				v-for="({ url, name }, j) in names"
				:key="j"
				@click="FetchAPI({ endpoint: url, parameters: '', _type: 'spellinfo' })">
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
