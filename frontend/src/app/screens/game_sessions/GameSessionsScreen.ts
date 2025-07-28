import { Container } from "pixi.js";
import { ConnectSocket } from "./websocket";

export class GameSessionsScreen extends Container {
  public static assetBundles = ["game_sessions"];

  constructor() {
    super();
  }
}
export class XoScreen extends Container {
  public static assetBundles = ["xo_game"];
  // private readonly bgColor = new Color(0x0e756b);

  constructor() {
    super();
    // ChangeBackgroundColorAnimation(this.bgColor, 1000);
    ConnectSocket("/api/game_sessions/");
  }

  public resize(width: number, height: number) {
    console.log(width, height);
    // const base = Math.min(width, height);
    // const cellSize = base * (this.boardScale / this.board.gridSize);

    // const emptyScale = 1 - this.boardScale;
  }
}
