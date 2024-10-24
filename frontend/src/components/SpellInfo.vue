<script setup lang="ts">
import type { SpellInfo } from "@/types/api";

const props = defineProps<{
	spell_info: SpellInfo | null;
}>();
const info = props.spell_info;
</script>

<!-- TODO: Replace names like Range, Duration... to icons -->
<template v-if="info != null">
	<h1>[Level {{ info!.level }}] {{ info!.name }}</h1>
	<p class="sub_info">Casting time {{ info!.casting_time }}</p>
	<p class="sub_info">Duration {{ info!.duration }}</p>
	<p class="sub_info">Range {{ info!.range }}</p>
	<p v-if="info!.material" class="sub_info">Material {{ info!.material }}</p>
	<p v-if="info!.concentration" class="sub_info">Triggers Concentration</p>
	<p v-for="(desc, i) in info!.desc" :key="i">{{ desc }}</p>
</template>

<style scoped lang="scss">
.sub_info {
	font-size: 12px;
	font-style: italic;
}

p:not(.sub_info) {
	font-size: 14px;
}
</style>
