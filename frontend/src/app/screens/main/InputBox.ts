import { Container, Point } from "pixi.js";
import { GetEngine } from "../../getEngine";

function DetectDirection(text: string): "rtl" | "ltr" {
	const rtlChars = /[\u0591-\u07FF\uFB1D-\uFDFD\uFE70-\uFEFC]/;
	return rtlChars.test(text) ? "rtl" : "ltr";
}

export class InputBox extends Container {
	private sizeRatio: Point;
	private htmlArea: HTMLTextAreaElement;


	constructor(
		placeholder = "Enter text...",
		sizeRatio: Point = new Point(0.8, 0.06),
	) {
		super();

		this.sizeRatio = sizeRatio;

		this.htmlArea = document.createElement("textarea");
		this.htmlArea.placeholder = placeholder;

		const style = this.htmlArea.style;
		style.position = "absolute";
		style.fontSize = "24px";
		style.padding = "8px";
		style.borderRadius = "8px";
		style.border = "none";
		style.outline = "none";
		style.background = "#333";
		style.color = "#fff";
		style.zIndex = "1000";
		style.resize = "none"; // prevent manual resizing
		style.overflowY = "auto"; // allow scroll when text overflows
		style.whiteSpace = "pre-wrap"; // preserve newlines
		style.lineHeight = "1.4";
		style.direction = DetectDirection(this.htmlArea.placeholder);

		this.htmlArea.oninput = () => {
			let value = this.htmlArea.value.trim();
			if (value === "") {
				value = this.htmlArea.placeholder;
			}

			const direction = DetectDirection(value);
			this.htmlArea.dir = direction;
			this.htmlArea.style.textAlign = direction === "rtl" ? "right" : "left";

			// 🔽 Auto-resize textarea height
			// this.htmlArea.style.height = "auto";
			// this.htmlArea.style.height = `${this.htmlArea.scrollHeight}px`;
		};

		document.body.appendChild(this.htmlArea);
	}

	public resize(width: number, height: number) {
		const canvasBounds = GetEngine().canvas.getBoundingClientRect();
		const boxWidth = canvasBounds.width * this.sizeRatio.x;
		const boxHeight = canvasBounds.height * this.sizeRatio.y;

		const yOffset = (canvasBounds.width - boxWidth) / 2;

		const style = this.htmlArea.style;
		style.width = `${boxWidth}px`;
		style.height = `${boxHeight}px`;
		style.left = `${yOffset}px`;
		style.right = `${yOffset}px`;
		style.bottom = `${boxHeight / 2}px`;
	}

	public destroy(options?: any) {
		this.htmlArea.remove();
		super.destroy(options);
	}
}
