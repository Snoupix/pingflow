export type Error = {
	error: string;
};

export type API_OUT_T = "classes" | "class" | "spells";
export type API_OUT = Classes | Class | SpellResp;

export type Classes = {
	count: number;
	results: Array<{ index: string; name: string; url: string }>;
};

// Generated using https://app.quicktype.io/?l=ts
export type Class = {
	index: string;
	name: string;
	hit_die: number;
	proficiency_choices: ProficiencyChoice[];
	proficiencies: Proficiency[];
	saving_throws: Proficiency[];
	starting_equipment: StartingEquipment[];
	starting_equipment_options: StartingEquipmentOption[];
	class_levels: string;
	multi_classing: MultiClassing;
	subclasses: Proficiency[];
	spellcasting?: Spellcasting;
	spells: string;
	url: string;
};

export type SpellResp = {
	count: number;
	results: Array<Spell>;
};

export type Spell = {
	index: string;
	name: string;
	level: number;
	url: string;
};

export type MultiClassing = {
	prerequisites: Prerequisite[];
	proficiencies: Proficiency[];
	proficiency_choices: ProficiencyChoice[];
};

export type Prerequisite = {
	ability_score: Proficiency;
	minimum_score: number;
};

export type Proficiency = {
	index: string;
	name: string;
	url: string;
};

export type ProficiencyChoice = {
	choose: number;
	type: string;
	from: ProficiencyChoiceFrom;
	desc?: string;
};

export type ProficiencyChoiceFrom = {
	option_set_type: string;
	options: PurpleOption[];
};

export type PurpleOption = {
	option_type: OptionType;
	item: Proficiency;
};

export enum OptionType {
	Reference = "reference",
}

export type Spellcasting = {
	level: number;
	spellcasting_ability: Proficiency;
	info: Info[];
};

export type Info = {
	name: string;
	desc: string[];
};

export type StartingEquipment = {
	equipment: Proficiency;
	quantity: number;
};

export type StartingEquipmentOption = {
	desc: string;
	choose: number;
	type: string;
	from: StartingEquipmentOptionFrom;
};

export type StartingEquipmentOptionFrom = {
	option_set_type: string;
	options: FluffyOption[];
};

export type FluffyOption = {
	option_type: string;
	count?: number;
	of?: Proficiency;
	choice?: Choice;
};

export type Choice = {
	desc: string;
	choose: number;
	type: string;
	from: ChoiceFrom;
};

export type ChoiceFrom = {
	option_set_type: string;
	equipment_category: Proficiency;
};
