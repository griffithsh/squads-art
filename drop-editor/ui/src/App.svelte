<script lang="ts">
  import Drop from "./lib/Drop.svelte";
  import type { Drop as DropType } from "src/types/Drop";

  const drop: DropType = {
    kind: "drop",
    key: "ABCCDF",
    description: "a description of this test drop",

    chances: [
      { kind: "nothing", probability: 3 },
      { kind: "item", code: "AXE", probability: 1 },
      {
        kind: "drop",
        key: "GGJDKRK",
        description: "dummy",
        chances: [{ kind: "nothing", probability: 1 }],
        probability: 1,
      },
    ],
  };
</script>

<main>
  <Drop {drop} />
</main>

<style>
  :root {
    font-family: Cantarell, "Bitstream Vera", DejaVu;
  }

  main {
    text-align: center;
    padding: 1em;
    margin: 0 auto;
  }
</style>
