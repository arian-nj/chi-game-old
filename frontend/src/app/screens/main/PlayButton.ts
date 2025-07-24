import { Container, FederatedPointerEvent, Sprite, Text } from "pixi.js";
import { gsap } from "gsap";
import { Signal } from "../../../engine/signals/signals";

const defaultScale = 1;
const bigScale = 1.1;

export class PlayButton extends Container {
  private bgSprite;
  private text;
  public OnClick = new Signal<void>();

  constructor() {
    super();
    // play button bg
    this.bgSprite = Sprite.from("play_btn_bg.png");
    this.bgSprite.anchor = 0.5;
    this.bgSprite.scale = defaultScale;
    this.addChild(this.bgSprite);
    this.eventMode = "dynamic";
    this.interactive = true;
    this.idleAnimation();

    this.onpointerdown = (_event: FederatedPointerEvent) => {
      this.stopAnimation();
    };

    this.onpointerup = (_event: FederatedPointerEvent) => {
      this.bounceAnimation();
      this.OnClick.Emit();
    };
    this.onpointerout = (_event: FederatedPointerEvent) => {
      this.continueAnimation();
    };

    // play button text
    this.text = new Text({
      text: "Play",
      style: {
        fill: "#ffffff",
        fontSize: 70,
        fontWeight: "900",
        // fontFamily: 'MyFont',
      },
      anchor: 0.5,
    });
    this.addChild(this.text);
  }
  private stopAnimation() {
    gsap.getTweensOf(this).forEach((t) => t.pause());
    // gsap.to(this, { pixi: { scale: bigScale }, duration: 0.7 })
  }

  private idleAnimation() {
    // y continuaslly grows and shrinks
    gsap.to(this, {
      pixi: { scaleX: 1.05, scaleY: 1.02 },
      duration: 1,
      ease: "power1.inOut",
      yoyo: true,
      repeat: -1,
    });
  }

  private continueAnimation() {
    gsap.getTweensOf(this).forEach((t) => t.play());
  }
  private bounceAnimation() {
    // gsap.to(this, {
    // 	pixi: { scale: 1.2 },
    // 	ease: "bounce.out",
    // 	duration: 0.7,
    // }).then(() =>
    // 	gsap.to(this, {
    // 		pixi: { scale: defaultScale },
    // 		duration: 0.7,
    // 		ease: "bounce.out",
    // 	})).then(() => {
    // 		this.idleAnimation()
    // 	})
    const tl = gsap.timeline();
    tl.to(this, { pixi: { scale: bigScale }, duration: 0.2 }).to(this, {
      pixi: { scale: defaultScale },
      duration: 0.3,
      ease: "bounce.out",
    });
  }
}
