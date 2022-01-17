<script lang="ts">
  import type { Probability } from "src/types/Drop";

  export let probability: Probability;
</script>

{probability.percent.toPrecision(5)}%: {probability.chance.kind}
{#if probability.children}
  <table>
    <tr>
      {#each probability.children as child}
        <td width="{probability.percent}%">
          <svelte:self probability={child} />
        </td>
      {/each}
    </tr>
  </table>
{/if}
