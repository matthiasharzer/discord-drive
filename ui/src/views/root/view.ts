import '../../components/Spinner';
import { css, html } from 'lit';
import { state } from 'lit/decorators.js';
import { Component } from '../../litutil/Component';

export class RootView extends Component {
	static styles = css`
		:host {
		}

		.view {
			display: flex;
			flex-direction: column;
			align-items: center;
			justify-content: center;

		}

		.layout {
			display: flex;
			flex-direction: column;
			gap: 1rem;
			max-width: 800px;
		}

		.body {
			display: flex;
			flex-direction: column;
			align-items: flex-start;
			gap: 0.5rem;
		}

		.neo-element {
			padding: 0.5rem;
		}

		.preview {
			padding: 0.1rem 0.5rem;
			font-weight: bold;
		}

		.actions {
			width: 100%;
			display: flex;
			flex-direction: row;
			justify-content: space-between;
			gap: 0.5rem;
		}

		.upload {
			display: flex;
			flex-direction: row;
			align-items: center;
			gap: 0.5rem;
		}

		h1 {
			font-size: 3rem;
			font-weight: bold;
			width: auto;
		}

		.select-file input {
			display: none;
		}

		.success {

			.content {
				display: flex;
				flex-direction: row;
				align-items: center;
				gap: 0.5rem;
				width: 100%;

				.copy {
					flex: 0 0 auto;
				}

				.download {
					flex: 1 1 auto;
					overflow: hidden;
					text-overflow: ellipsis;
					white-space: nowrap;
				}
			}
		}

		svg {
			height: 1em;
			width: 1em;
		}

		a {
			color: var(--primary-color, #fff);
			font-weight: bold;
		}
	`;

	@state()
	selectedFile: File | null = null;

	@state()
	uploading: boolean = false;

	@state()
	errorMessage: string | null = null;

	@state()
	downloadURL: string | null = null;

	@state()
	copied: boolean = false;

	selectFile(event: Event) {
		const input = event.target as HTMLInputElement;
		if (input.files && input.files.length > 0) {
			this.selectedFile = input.files[0];
		} else {
			this.selectedFile = null;
		}
	}

	uploadFile() {
		if (!this.selectedFile) {
			return;
		}
		if (this.uploading) {
			return;
		}

		this.errorMessage = null;
		this.downloadURL = null;
		this.uploading = true;

		const encodedFileName = encodeURIComponent(this.selectedFile.name);

		fetch(`/api/v1/upload/${encodedFileName}`, {
			method: 'POST',
			body: this.selectedFile,
		})
			.then(response => {
				if (!response.ok) {
					throw new Error(`Upload failed with status ${response.status}`);
				}
				return response;
			})
			.then(() => {
				this.selectedFile = null;
				this.downloadURL = `${window.location.origin}/api/v1/file/${encodedFileName}`;
			})
			.catch(error => {
				this.errorMessage = `Error uploading file: ${error.message}`;
				console.error('Error uploading file:', error);
			})
			.finally(() => {
				this.uploading = false;
			});
	}

	copyLink() {
		if (this.downloadURL) {
			navigator.clipboard
				.writeText(this.downloadURL)
				.then(() => {
					console.log('Link copied to clipboard');
					this.copied = true;
					setTimeout(() => {
						this.copied = false;
					}, 2000);
				})
				.catch(err => {
					console.error('Could not copy link: ', err);
				});
		}
	}

	renderSuccess() {
		if (!this.downloadURL) {
			return html``;
		}
		return html`
			<div class="neo-element body thick success">
				<h2 class="neo-element header thin">File Uploaded Successfully</h2>
				<div class="content">
					<button class="neo-element thin secondary copy" @click=${this.copyLink}>
						${
							this.copied
								? html`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" height="24px" width="24px" fill="#e3e3e3"><path d="M9.55,18.8l-6.5-6.5L5.3,10.05L9.55,14.3L18.7,5.15L20.95,7.4L9.55,18.8Z"/></svg>`
								: html`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" height="24px" width="24px" fill="#e3e3e3"><path d="M6.73,18.27V0.57h14.7v17.7H6.73Zm3.15-3.15h8.4V3.72H9.88v11.4Zm-8.3,8.3V5.72H4.73V20.27H16.28v3.15H1.58Zm8.3-8.3V3.72v11.4Z"/></svg>`
						}
					</button>
					<div class="neo-element header thin download">
						<a href="${this.downloadURL}" target="_blank">${this.downloadURL}</a>
					</div>
				</div>
			</div>
		`;
	}

	render() {
		return html`
			<div class="view">
				<div class="layout">
					<div class="neo-element header thick">
						<h1>Discord Drive</h1>
					</div>
					<div class="neo-element body thick">
						<div class="actions">
							<label class="neo-element thin select-file button">
									Select File
									<input type="file" @change=${this.selectFile}  />
							</label>

							<button class="neo-element thin secondary upload" ?disabled=${!this.selectedFile || this.uploading} @click=${this.uploadFile}>
								${
									this.uploading
										? html`<dd-spinner></dd-spinner>`
										: html`
								<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" height="24px" width="24px" fill="#e3e3e3"><path d="M10.43,20.85V9.2L5.38,14.25L3.15,12L12,3.15L20.85,12l-2.23,2.25L13.58,9.2V20.85H10.43Z"/></svg>`
								}
								Upload
							</button>
						</div>
							${
								this.selectedFile
									? html`
									<div class="neo-element header thin preview">
										<span>${this.selectedFile.name}</span>
									</div>`
									: ''
							}
							${
								this.errorMessage
									? html`
									<div class="neo-element header thin error">
										<span>${this.errorMessage}</span>
									</div>`
									: ''
							}
					</div>
					${this.downloadURL ? this.renderSuccess() : ''}
				</div>
			</div>
		`;
	}
}

customElements.define('dd-root-view', RootView);
