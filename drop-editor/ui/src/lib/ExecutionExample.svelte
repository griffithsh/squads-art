<script lang="ts">
  import mt19937 from "@stdlib/random-base-mt19937";
  import type { Drop, Item, Nothing } from "src/types/Drop";

  export let drop: Drop;
  export let seed: number;

  const prng = mt19937.factory({
    seed,
  });

  const recurse = (drop: Drop): Item | Nothing => {
    let sum = drop.chances.reduce((total, chance) => {
      return total + chance.probability;
    }, 0);

    const roll = prng() % sum;
    let running = 0;
    let rolled: Item;
    for (let i = 0; i < drop.chances.length; i++) {
      if (roll < drop.chances[i].probability + running) {
        if (drop.chances[i].kind === "drop") {
          return recurse(drop.chances[i] as Drop);
        }
        return drop.chances[i] as Item | Nothing;
      }
    }
  };

  const result = recurse(drop);
</script>

<p>
  This could be what you get.
  {#if result.kind === "item"}
    Item! {result.code}
  {:else}
    Unlucky! No drop.
  {/if}
</p>
