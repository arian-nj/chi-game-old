import { Color, Container } from "pixi.js";
import { ChangeBackgroundColorAnimation } from "../../animateion/change_background";
import { XoBoard } from "./XOBoard";
export class XoScreen extends Container {
  public static assetBundles = ["xo_game"];
  private readonly bgColor = new Color(0x0e756b);
  // private readonly gridSize = 3;
  private boardScale = 0.9;

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

    const emptyScale = 1 - this.boardScale;
    this.board.position.set(
      (width * emptyScale) / 2,
      (height * emptyScale) / 2,
    );
    this.board.resize(cellSize);
  }
}
