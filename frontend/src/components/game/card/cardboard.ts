import Phaser from "phaser";
import { Card } from "./card";


class CardGame extends Phaser.Scene {
  constructor() {
    super({ key: 'CardGame' });
  }

  preload(): void {
    this.load.image('cardBack', 'assets/card/CardBack.png');

    // const cardValues: string[] = ['A', 'K', 'Q', 'J', '10', '9'];
    this.load.image('bgCard', 'assets/card/BgCard.png');
  }

  create(): void {
    this.cameras.main.setBackgroundColor('#2E8B57');

    const cardValues: string[] = ['A', 'K', 'Q', 'J', '10', '9'];

    // Position and create cards.
    const startX: number = 150;
    const startY: number = 300;
    const spacing: number = 160;

    cardValues.forEach((value, index) => {
      // Create a new Card instance.
      new Card(this, startX + (index * spacing), startY, value);
    });

  }
}

export function CreateGame(parentId: string) {
  return new Phaser.Game({
    type: Phaser.AUTO,
    width: 800,
    height: 600,
    backgroundColor: "#1d1d1d",
    parent: parentId,
    scene: [CardGame]
  });
}
