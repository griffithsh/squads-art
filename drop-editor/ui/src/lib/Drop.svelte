<script lang="ts">
  import type { Drop, Item } from "src/types/Drop";
import ExecutionExample from "./ExecutionExample.svelte";
  import Probabilities from "./Probabilities.svelte";

  export let drop: Drop;

  const probabilities = (drop: Drop) => {
    let sum = drop.chances.reduce((total, chance) => {
      return total + chance.probability;
    }, 0);

    return drop.chances.map((chance) => {
      return {
        percent: (chance.probability / sum) * 100,
        chance,
        children: chance.kind === "drop" ? probabilities(chance) : undefined,
      };
    });
  };
</script>

<p>What drops?</p>
<div>{drop.key} - {drop.description}</div>
<ul>
  {#each drop.chances ?? [] as chance}
    <li>
      {#if chance.kind === "nothing"}
        Nothing
      {:else if chance.kind === "drop"}
        Continue with: {chance.key}
      {:else if chance.kind === "item"}
        Got item: {chance.code}
      {/if}
    </li>
  {/each}
</ul>
<div>
  <Probabilities {drop} />
</div>
<div>
  <ExecutionExample seed={12345} {drop}/>
</div>

<style>
</style>
