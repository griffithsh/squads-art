<script lang="ts">
  import type { Drop, Probability } from "src/types/Drop";
  import Row from "./Probability/Row.svelte";

  export let drop: Drop;

  const probabilities = (drop: Drop): Probability[] => {
    let sum = drop.chances.reduce((total, chance) => {
      return total + chance.probability;
    }, 0);

    return drop.chances.map((chance): Probability => {
      return {
        percent: (chance.probability / sum) * 100,
        chance,
        children: chance.kind === "drop" ? probabilities(chance) : undefined,
      };
    });
  };
</script>

<table>
  <tr>
    {#each probabilities(drop) as probability}
      <td width="{probability.percent}%">
        <Row {probability} />
      </td>
    {/each}
  </tr>
</table>
