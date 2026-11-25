import { type CSSResultGroup, css, LitElement } from 'lit';

const cssOverwrites = css`

input {
  min-width: 0;
  width: auto;
	font-family: inherit;
}

*, *::before, *::after {
  box-sizing: border-box;
}
* {
  margin: 0;
	interpolate-size: allow-keywords;
	font-variant-ligatures: none;
	font-family: var(--font-body);
}

img, picture, video, canvas, svg {
  display: block;
  max-width: 100%;
}
input, button, textarea, select {
  font: inherit;
}
p, h1, h2, h3, h4, h5, h6 {
  overflow-wrap: break-word;
}
p {
  text-wrap: pretty;
}
h1, h2, h3, h4, h5, h6 {
  text-wrap: balance;
	font-family: var(--font-header);
}
#root, #__next {
  isolation: isolate;
}

.neo-element {
	--background-res: var(--background, #a9d7b8);
	--ink-res: var(--ink, #000);
	--border-width-res: var(--border-width, 3px);
	--border-color-res: var(--border-color, #000);
	--shadow-offset-x-res: var(--shadow-offset-x, 4px);
	--shadow-offset-y-res: var(--shadow-offset-y, 4px);
	--shadow-color-res: var(--shadow-color, #000);

	background-color: var(--background-res);
	color: var(--ink-res);
	border: var(--border-width-res) solid var(--border-color-res);
	box-shadow: var(--shadow-offset-x-res) var(--shadow-offset-y-res) 0px var(--shadow-color-res);

	&.header {
		--background: var(--header-background-color, #FF90E8);
		--ink: var(--header-ink-color, #000);
	}

	&.body {
		--background: var(--body-background-color, #a9d7b8);
		--ink: var(--body-ink-color, #000);
	}

	&.thick {
		--shadow-offset-x: 8px;
		--shadow-offset-y: 8px;
		--border-width: 6px;
	}

	&.thin {
		--shadow-offset-x: 4px;
		--shadow-offset-y: 4px;
		--border-width: 3px;
	}

	&.secondary {
		--background: var(--secondary-background-color, #5662f6);
		--ink: var(--secondary-ink-color, #fff);
	}
}

button.neo-element, .button.neo-element {
	--background: var(--button-background-color, #5662f6);
	--ink: var(--button-ink-color, #fff);

	cursor: pointer;
	transition: all 0.1s ease-in-out;

	&:hover {
		transform: translate(-2px, -2px);
		box-shadow: calc(var(--shadow-offset-x, 4px) + 2px) calc(var(--shadow-offset-y, 4px) + 2px) 0px var(--shadow-color, #000);
	}

	&:active {
		transform: translate(calc(var(--shadow-offset-x-res)),calc(var(--shadow-offset-y-res)));
		box-shadow: 0 0 0px var(--shadow-color, #000);
	}

	&:disabled {
		filter: grayscale(100%);
		pointer-events: none;
	}
}
`;

export class Component extends LitElement {
	protected static _styles: CSSResultGroup;

	static get styles(): CSSResultGroup {
		const derivedStyles = Component._styles || [];
		return [cssOverwrites, ...(Array.isArray(derivedStyles) ? derivedStyles : [derivedStyles])];
	}

	static set styles(styles: CSSResultGroup) {
		Component._styles = styles;
	}

	get rect() {
		return this.getBoundingClientRect();
	}

	sleep(ms: number) {
		return new Promise(resolve => setTimeout(resolve, ms));
	}

	connectedCallback(): void {
		super.connectedCallback();
	}

	dispatch<T>(
		name: string,
		detail: T | null = null,
		options: { bubbles?: boolean; composed?: boolean } = {},
	) {
		this.dispatchEvent(new CustomEvent(name, { detail, ...options }));
	}

	dispatchBubble<T>(name: string, detail: T | null = null) {
		this.dispatch(name, detail, { bubbles: true, composed: true });
	}

	disconnectedCallback() {
		super.disconnectedCallback();
	}
}
