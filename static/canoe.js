CanoeFragment = window.CanoeFragment || class extends HTMLElement {
  static get is() { return "canoe-fragment" }

  set fragment(val) {
    this.setImport(val)
  }

  setImport(val) {
    let link = document.createElement('link')
    link.rel = "import"
    link.href = "/fragments/" + val
    link.setAttribute('async', '')
    link.onload =  this.importLoaded
    this.innerHTML = ''
    this.appendChild(link)
  }

  importLoaded() {
    parent = this.parentNode
    const t = this.import.querySelector("template")
    if(t !== null) {
      const instance = t.content.cloneNode(true);
      parent.shadowRoot.appendChild(instance);
    }
  }

  connectedCallback() {
    this.setImport(this.frag)
  }

  constructor() {
    super();
    this.attachShadow({mode: 'open'});
    this.frag = this.getAttribute("fragment")
  }
}

if (typeof window.customElements.get(CanoeFragment.is) === "undefined") {
  window.customElements.define(CanoeFragment.is, CanoeFragment);
}
