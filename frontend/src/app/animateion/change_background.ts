import { Color, Ticker } from "pixi.js";
import { GetEngine } from "../getEngine";
import { lerp } from "../../engine/utils/maths";


export function ChangeBackgroundColorAnimation(endColor: Color, duration: number) {
	const eng = GetEngine();

	let elapsed = 0;

	const startColor = eng.renderer.background.color.toRgb(); // red
	const endColorRgb = endColor.toRgb();

	const animation = (ticker: Ticker) => {
		elapsed += ticker.elapsedMS;
		const t = Math.min(elapsed / duration, 1);
		const r = lerp(startColor.r, endColorRgb.r, t) * 255;
		const g = lerp(startColor.g, endColorRgb.g, t) * 255;
		const b = lerp(startColor.b, endColorRgb.b, t) * 255;

		const new_color = new Color({ r: r, g: g, b: b });

		eng.renderer.background.color.setValue(new_color);

		if (t >= 1) {
			eng.ticker.remove(animation); // optional: stop once done
		}
	};
	eng.ticker.add(animation)
}
