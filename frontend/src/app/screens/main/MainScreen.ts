import { Container } from "pixi.js";
import { PlayButton } from "./PlayButton";
import { GetJwtToken } from "../../../auth/jwt";

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

  public async onClick() {
    console.log("clicked from menu");
    ConnectMatchmakingSocket();
  }
}
function ConnectMatchmakingSocket() {
  const url = "/api/game/match_making/ticket?auth_token=" + GetJwtToken();
  const socket = new WebSocket(url, []);
  console.log(socket);

  socket.onopen = function (_e: Event) {
    alert("[open] Connection established");
    console.log(_e);
    socket.send("My name is John");
  };

  socket.onmessage = function (event) {
    alert(`[message] Data received from server: ${event.data}`);
    socket.close(1000, "Normal Closure");
  };

  socket.onclose = function (event) {
    if (event.wasClean) {
      alert(
        `[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`,
      );
    } else {
      alert("[close] Connection died");
    }
  };
}
