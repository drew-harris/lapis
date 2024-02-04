import "vite/modulepreload-polyfill";
import { render } from "solid-js/web";
import { StandardsForm } from "./StandardsForm";

const islands: Record<string, () => any> = {
  standards: StandardsForm,
};

class Island extends HTMLElement {
  rendered = false;
  dispose?: () => void;
  async connectedCallback() {
    const src = this.getAttribute("src") ?? "";
    const componentLoader = islands[src];
    if (!componentLoader) {
      throw new Error(`${src} is not a component! Check your islands/index.`);
    }

    this.dispose = render(componentLoader, this);
    this.rendered = true;
  }
}

customElements.define("vite-land", Island);

for (let element of document.querySelectorAll("vite-land")) {
  const src = element.getAttribute("src") ?? "";
  const componentLoader = islands[src];
  if (!componentLoader) {
    throw new Error(`${src} is not a component! Check your islands/index.`);
  }
  if (element instanceof Island) {
    if (!element.rendered) {
      element.dispose = render(componentLoader, element);
      element.rendered = true;
    }
  }
}
