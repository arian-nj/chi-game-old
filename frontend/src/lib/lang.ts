export function detectDirection(text: string): "rtl" | "ltr" {
	const firstStrongChar = text.match(/[\u0591-\u07FF\uFB1D-\uFDFD\uFE70-\uFEFC\w]/);
	if (!firstStrongChar) return "ltr"; // default
	const char = firstStrongChar[0];
	return /[\u0591-\u07FF\uFB1D-\uFDFD\uFE70-\uFEFC]/.test(char) ? "rtl" : "ltr";
}
