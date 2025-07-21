import gsap from "gsap";
import { Container, FederatedPointerEvent, Sprite } from "pixi.js";
import { SettingsPopup } from "../popups/SettingsPopup";
import { GetEngine } from "../getEngine";
const defaultScale = 2;
const bigScale = 2.5;

export class SettingButton extends Container {
  iconSprite: Sprite;
  constructor() {
    super();
    this.iconSprite = Sprite.from("icon-settings.png");
    this.iconSprite.scale = defaultScale;
    this.iconSprite.eventMode = "static";
    this.iconSprite.anchor = 0.5;

    this.iconSprite.onpointerdown = (_event: FederatedPointerEvent) => {
      gsap.to(this.iconSprite, { pixi: { scale: bigScale }, duration: 0.7 });
      console.log("down");
    };

    this.iconSprite.onpointerup = (_event: FederatedPointerEvent) => {
      gsap.to(this.iconSprite, {
        pixi: { rotation: "+=360_cw", scale: defaultScale },
        duration: 0.7,
      });
      GetEngine().navigation.presentPopup(SettingsPopup);
    };

    this.iconSprite.onpointerout = (_event: FederatedPointerEvent) => {
      gsap.to(this.iconSprite, {
        pixi: { rotation: "+=360_cw", scale: defaultScale },
        duration: 0.7,
        ease: "bounce.out",
      });
    };

    this.addChild(this.iconSprite);
  }

  public resize(width: number, _height: number) {
    this.x = width - this.width;
    this.y = this.height;
  }
}
