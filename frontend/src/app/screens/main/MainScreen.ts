import { Container } from "pixi.js";
import { PlayButton } from "./PlayButton";

export class MainScreen extends Container {
  public static assetBundles = ["menu"];
  private playButton;
  constructor() {
    super();

    this.playButton = new PlayButton();
    this.playButton.OnClick.Connect(this.onClick.bind(this));
    this.addChild(this.playButton);
  }

  public resize(width: number, height: number) {
    this.playButton.x = width / 2;
    this.playButton.y = height / 2;
  }

  public onClick() {
    console.log("clicked from menu");
  }
}
