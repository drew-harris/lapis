import { createSignal } from "solid-js";

export function StandardsForm() {
  const [name, setName] = createSignal("drew was here");
  return (
    <>
      <div>{name() + ""}</div>
      <div>{name() + ""}</div>
      <input
        class="bg-purple-200"
        type="text"
        value={name()}
        onInput={(e) => setName((e.target as HTMLInputElement).value)}
      />
    </>
  );
}
