import {
  Color,
  Container,
  FederatedPointerEvent,
  Graphics,
  Point,
  Rectangle,
  Ticker,
} from "pixi.js";
import { ChangeBackgroundColorAnimation } from "../../animateion/change_background";
import { GetEngine } from "../../getEngine";
import { lerp } from "../../../engine/utils/maths";
import { ConnectSocket } from "./xo_websocket";

enum CellType {
  EmptyCell = "",
  XCell = "X",
  OCell = "O",
}

export class XoScreen extends Container {
  public static assetBundles = ["xo_game"];
  private readonly bgColor = new Color(0x0e756b);
  // private readonly gridSize = 3;
  private boardScale = 0.9;

  private readonly board: XoBoard;

  constructor() {
    super();
    ConnectSocket("/api/game/xo/ws/");
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

class XoBoard extends Container {
  public gridSize = 3;
  private readonly lineColor = new Color(0x0da192);
  private cellSize = 0;
  private lineWidth = 20;

  private boardGraphic = new Graphics();
  private selectorGraphics = new Graphics();

  private cellsGraphics = new Graphics();
  private Cells: CellType[] = [];
  private currentPlayer: CellType = CellType.XCell;

  constructor(gridSize: number) {
    super();
    this.gridSize = gridSize;
    for (let i = 0; i <= gridSize * gridSize; i++) {
      this.Cells.push(CellType.EmptyCell);
    }

    this.eventMode = "static";
    this.addChild(this.boardGraphic);
    this.addChild(this.selectorGraphics);
    this.addChild(this.cellsGraphics);
    this.onmousemove = (event: FederatedPointerEvent) => {
      const cellCord = this.cellCordFromPos(event);
      this.selectorGraphics.clear();

      const x = cellCord.y * this.cellSize + this.lineWidth / 2;
      const y = cellCord.x * this.cellSize + this.lineWidth / 2;

      const colorCellSize = this.cellSize - this.lineWidth;
      this.selectorGraphics.rect(x, y, colorCellSize, colorCellSize);

      this.selectorGraphics.fill({ color: 0xff0000, alpha: 0.3 });
    };

    this.onmouseup = (event: FederatedPointerEvent) => {
      const cellCord = this.cellCordFromPos(event);
      const index = cellCord.x * this.gridSize + cellCord.y;

      // Skip if the cell is already taken
      if (this.Cells[index] !== CellType.EmptyCell) {
        return;
      }

      // Mark cell
      this.Cells[index] = this.currentPlayer;
      this.drawSymbol(cellCord.x, cellCord.y, this.currentPlayer);

      // Switch turn
      this.currentPlayer =
        this.currentPlayer === CellType.XCell ? CellType.OCell : CellType.XCell;
    };
  }

  private drawSymbol(row: number, col: number, symbol: CellType) {
    const x = col * this.cellSize + this.lineWidth / 2;
    const y = row * this.cellSize + this.lineWidth / 2;
    const size = this.cellSize - this.lineWidth;
    const engine = GetEngine();
    if (symbol === CellType.XCell) {
      // Draw X
      let elapse = 0;
      const duration = 1000;
      const animateX = (ticker: Ticker) => {
        this.cellsGraphics.clear();
        elapse += ticker.elapsedMS;
        const progress = Math.min(elapse / duration, 1);
        const posX = lerp(x, x + size, progress);
        this.cellsGraphics
          .moveTo(x, y)
          .lineTo(posX, y + this.cellSize)
          .stroke({ color: 0x8c1837, width: 20 });

        if (progress >= 1) {
          engine.ticker.remove(animateX);
        }
      };
      engine.ticker.add(animateX);
      // this.cellsGraphics.moveTo(x, y)
      // 	.lineTo(x + size, y + size)
      // 	.moveTo(x + size, y)
      // 	.lineTo(x, y + size)
      // 	.stroke({ color: 0x8c1837, width: 20 });
    } else if (symbol === CellType.OCell) {
      // Draw O
      const centerX = x + size / 2;
      const centerY = y + size / 2;
      const radius = size / 2;
      this.cellsGraphics
        .circle(centerX, centerY, radius)
        .stroke({ color: 0x1a3caa, width: 20 });
    }
  }

  private cellCordFromPos(event: FederatedPointerEvent): Point {
    const local = this.toLocal(event.global);
    const row = Math.floor(local.y / this.cellSize);
    const col = Math.floor(local.x / this.cellSize);
    return new Point(row, col);
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
