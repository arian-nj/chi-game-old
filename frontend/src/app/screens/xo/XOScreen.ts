import { Color, Container, FederatedPointerEvent, Graphics, Rectangle } from "pixi.js";
import { ChangeBackgroundColorAnimation } from "../../animateion/change_background";


export class XoScreen extends Container {
	public static assetBundles = ["xo_game"];
	private readonly bgColor = new Color(0x0e756b);
	// private readonly gridSize = 3;
	private boardScale = 0.9

	private readonly board: XoBoard;

	constructor() {
		super();
		ChangeBackgroundColorAnimation(this.bgColor, 1000);

		this.board = new XoBoard(3);
		this.addChild(this.board);
	}

	public resize(width: number, height: number) {
		const base = Math.min(width, height);
		const cellSize = base * (this.boardScale / this.board.gridSize);

		const emptyScale = 1 - this.boardScale
		this.board.position.set((width * emptyScale) / 2, (height * emptyScale) / 2);
		this.board.resize(cellSize);
	}
}
class XoBoard extends Container {
	public gridSize = 3;
	private readonly lineColor = new Color(0x0da192);
	private cellSize = 0;

	private boardGraphic = new Graphics()
	private selectorGraphics = new Graphics()
	private lineWidth = 20

	constructor(gridSize: number) {
		super();
		this.gridSize = gridSize
		this.eventMode = 'static';
		this.addChild(this.boardGraphic)
		this.addChild(this.selectorGraphics)
		this.onmousemove = (event: FederatedPointerEvent) => {
			const local = this.toLocal(event.global)
			const row = Math.floor(local.y / this.cellSize)
			const col = Math.floor(local.x / this.cellSize)
			console.log(row, col)
			this.selectorGraphics.clear()

			const x = col * this.cellSize + this.lineWidth / 2
			const y = row * this.cellSize + this.lineWidth / 2
			const colorCellSize = this.cellSize - this.lineWidth
			this.selectorGraphics.rect(x, y, colorCellSize, colorCellSize)

			this.selectorGraphics.fill({ color: 0xff0000, alpha: 0.3 });
		}
	}

	public resize(cellSize: number) {
		this.cellSize = cellSize;
		const boardSize = this.cellSize * this.gridSize;
		this.hitArea = new Rectangle(0, 0, boardSize, boardSize);
		this.drawBoard();
	}

	private drawBoard() {
		this.boardGraphic.clear();

		const boardSize = this.cellSize * this.gridSize;

		for (let i = 1; i < this.gridSize; i++) {
			const x = i * this.cellSize;
			this.boardGraphic.moveTo(x, 0);
			this.boardGraphic.lineTo(x, boardSize);
		}

		for (let i = 1; i < this.gridSize; i++) {
			const y = i * this.cellSize;
			this.boardGraphic.moveTo(0, y);
			this.boardGraphic.lineTo(boardSize, y);
		}

		this.boardGraphic.stroke({
			color: this.lineColor,
			pixelLine: false,
			width: this.lineWidth,
			alpha: 1,
			cap: "round",
		});
	}
}
