

export class Card extends Phaser.GameObjects.Sprite {
  value: string;
  isFlipped: boolean;
  isFlipping: boolean;

  constructor(scene: Phaser.Scene, x: number, y: number, value: string) {
    // The card starts with the 'cardBack' texture.
    super(scene, x, y, 'cardBack');

    this.scene = scene;
    this.value = value;
    this.setInteractive();
    this.scene.input.setDraggable(this)
    this.scene.add.existing(this);

    this.isFlipped = false;
    this.isFlipping = false;

    this.on('pointerdown', () => {
      this.flip();
    });

    this.on('drag', (_pointer: Phaser.Input.Pointer, dragX: number, dragY: number) => this.setPosition(dragX, dragY))
    this.on('dragstart', (_pointer: Phaser.Input.Pointer, _gameObject: Phaser.GameObjects.GameObject) => {
      this.scale = 1.0
    })
    this.on('dragend', (_pointer: Phaser.Input.Pointer, _gameObject: Phaser.GameObjects.GameObject) => {
      this.scale = 1
    })
  }

  flip(): void {
    if (this.isFlipping) {
      return;
    }
    this.isFlipping = true;

    this.scene.tweens.add({
      targets: this,
      scaleX: 0,
      ease: 'Linear',
      duration: 150,
      onComplete: () => {
        if (this.isFlipped) {
          this.setTexture('cardBack');
        } else {
          this.setTexture('bgCard');
        }
        this.isFlipped = !this.isFlipped;

        this.scene.tweens.add({
          targets: this,
          scaleX: 1,
          ease: 'Linear',
          duration: 150,
          onComplete: () => {
            this.isFlipping = false;
          }
        });
      }
    });
  }

}
