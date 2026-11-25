import './components/Background.ts';
import { Router } from '@lit-labs/router';
import { css, html } from 'lit';
import { Component } from './litutil/Component.ts';
import 'urlpattern-polyfill';

export class App extends Component {
	static styles = css`
		:host {
			display: flex;
			flex-direction: column;
			overflow-y: auto;
			scrollbar-gutter: stable both-edges;

			width: 100%;
			height: 100%;
			min-height: 100dvh;
		}

		main {
			flex: 1;
			display: flex;
			flex-direction: column;
			justify-content: safe center;

			width: 100%;
		}

		dd-background {
			z-index: -1;
		}
	`;

	private router = new Router(this, [
		{
			path: '/',
			render: () => html`<dd-root-view></dd-root-view>`,
			enter: async () => {
				await import('./views/root/view.ts');
				return true;
			},
		},
	]);

	render() {
		return html`
			<dd-background></dd-background>
			<main>
				${this.router.outlet()}
			</main>
		`;
	}
}

customElements.define('dd-app', App);
