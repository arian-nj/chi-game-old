import { Application, Container, Ticker } from "pixi.js";

export class Manager {
  public static app: Application;
  private static currentScene: IScene;

  public static get width(): number {
    return Math.max(
      document.documentElement.clientWidth,
      window.innerWidth || 0,
    );
  }

  public static get height(): number {
    return Math.max(
      document.documentElement.clientHeight,
      window.innerHeight || 0,
    );
  }

  public static async initialize(background: number): Promise<void> {
    Manager.app = new Application();

    await Manager.app.init({
      resizeTo: window,
      resolution: window.devicePixelRatio || 1,
      autoDensity: true,
      backgroundColor: background,
    });

    document.getElementById("pixi-container")!.appendChild(Manager.app.canvas);

    Manager.app.ticker.add(Manager.update);
    window.addEventListener("resize", Manager.resize);
  }
  public static resize(): void {
    if (Manager.currentScene) {
      Manager.currentScene.resize(Manager.width, Manager.height);
    }
  }

  public static changeScene(newScene: IScene): void {
    if (Manager.currentScene) {
      Manager.app.stage.removeChild(Manager.currentScene);
      Manager.currentScene.destroy();
    }

    Manager.currentScene = newScene;
    Manager.app.stage.addChild(Manager.currentScene);
  }

  private static update(ticker: Ticker): void {
    if (Manager.currentScene) {
      Manager.currentScene.update(ticker);
    }
  }
}

export interface IScene extends Container {
  update(ticker: Ticker): void;
  resize(width: number, height: number): void;
}
