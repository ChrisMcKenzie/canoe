(function(window) {
  window.CanoeFragment = window.CanoeFragment || class extends HTMLElement {
    static get is() { return "canoe-fragment" }

    set fragment(val) {
      this.setImport(val)
    }

    set loaded(val) {
      if(val) {
        this.setAttribute("loaded", "")
      } else {
        this.removeAttribute("loaded")
      }
    }

    set failed(val) {
      if(val) {
        this.setAttribute("failed", "")
      } else {
        this.removeAttribute("failed")
      }
    }

    setImport(val) {
      this.loaded = false
      this.failed = false
      let link = document.createElement('link')
      link.rel = "import"
      link.href = "/fragments/" + val
      // link.setAttribute('async', '')
      link.onload =  this.importLoaded
      link.onerror = this.handleError
      this.innerHTML = ''
      this.appendChild(link)
    }

    handleError(e) {
      this.parentNode.failed = true
    }

    importLoaded(e) {
      parent = this.parentNode
      const t = this.import.querySelector("template")
      if(t !== null) {
        const instance = t.content.cloneNode(true)
        parent.shadowRoot.appendChild(instance)
      }
      parent.loaded = true
      parent.failed = false
    }

    connectedCallback() {
      this.setImport(this.frag)
    }

    constructor() {
      super();
      this.attachShadow({mode: 'open'})
      this.frag = this.getAttribute("fragment")
    }
  }

  if (typeof window.customElements.get(CanoeFragment.is) === "undefined") {
    window.customElements.define(CanoeFragment.is, CanoeFragment)
  }
})(window);
