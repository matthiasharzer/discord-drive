import { css, html } from 'lit';
import { Component } from '../../litutil/Component';

export class RootView extends Component {
	static styles = css`
		:host {
		}
	`;

	render() {
		return html`
			<p>RootView</p>
		`;
	}
}

customElements.define('dd-root-view', RootView);
