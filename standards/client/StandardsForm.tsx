import { createSignal } from "solid-js";

export function StandardsForm() {
  const [name, setName] = createSignal("james");
  return (
    <>
      <div>{name() + ""}</div>
      <div>{name() + ""}</div>
      <input
        class="bg-green-100"
        type="text"
        value={name()}
        onInput={(e) => setName((e.target as HTMLInputElement).value)}
      />
    </>
  );
}
