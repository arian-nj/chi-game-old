import { Container, Ticker } from "pixi.js";
import { IScene } from "../Manager";

export class GameScene extends Container implements IScene {
  constructor() {
    super();
  }

  public update(ticker: Ticker): void {
    void ticker;
  }

  public resize(width: number, height: number): void {
    void width;
    void height;
  }
}
