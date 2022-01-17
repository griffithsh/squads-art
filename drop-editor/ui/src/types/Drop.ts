/**
 * There are entry points to the drop table that game code can reference. An
 * entry point should define how many items could be rolled as well as which
 * items could be rolled for. I think this is called a "Drop".
 *
 * They have N children which could each be one of a reference to some other
 * entry point, a concrete item, or nothing. Each child must be marked with a
 * probability of being resolved.
 */


export interface Nothing {
  kind: "nothing"
  probability: number
}
export interface Item {
  kind: "item";
  code: string;
  probability: number
}

export interface DropProbability extends Drop {
  probability: number
}

export interface Drop {
  kind: "drop";
  key: string;
  description: string;
  chances: (DropProbability | Item | Nothing)[];
}

export interface Probability {
  percent: number;
  chance: DropProbability | Item | Nothing;
  children: Probability[];
}
