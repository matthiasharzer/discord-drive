import { css, html } from 'lit';
import { Component } from '../litutil/Component.ts';

export class Background extends Component {
	static styles = css`
		:host {
			position: absolute;
			background-color: #fff6e5;
			top: 0;
			left: 0;
			right: 0;
			bottom: 0;
		}

		.dotted-overlay {
			position: absolute;
			top: 0;
			left: 0;
			right: 0;
			bottom: 0;

			 --circle-diameter: 2px;
			 --circle-spacing: 40px;
			 --circle-color: color-mix(in srgb, #000000 50%, transparent);
			background : radial-gradient(
				circle at
						var(--circle-diameter)
						var(--circle-diameter),
				var(--circle-color) calc(var(--circle-diameter) - 1px),
				transparent var(--circle-diameter)
			)
			0 0 / var(--circle-spacing) var(--circle-spacing);
		}
	`;

	render() {
		return html`
		<div class="dotted-overlay"></div>
		`;
	}
}

customElements.define('dd-background', Background);
