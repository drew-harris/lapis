import "vite/modulepreload-polyfill";
import { render } from "solid-js/web";
import { StandardsForm } from "./StandardsForm";

render(() => <StandardsForm />, document.getElementById("standards-form")!);
