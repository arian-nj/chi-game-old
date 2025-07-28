import { Container } from "pixi.js";
import { PlayButton } from "./PlayButton";
import { GetJwtToken } from "../../../auth/jwt";
import { GetEngine } from "../../getEngine";
import { GameSessionsScreen } from "../game_sessions/GameSessionsScreen";
import { InputBox } from "./InputBox";

export class MainScreen extends Container {
	public static assetBundles = ["menu"];
	private playButton;
	private inputBox = new InputBox("به حریفت پیام بده");

	constructor() {
		super();

		this.playButton = new PlayButton();
		this.playButton.OnClick.Connect(this.onClick.bind(this));
		this.addChild(this.playButton);

		this.addChild(this.inputBox);
	}

	public resize(width: number, height: number) {
		this.inputBox.resize(width, height);
		this.playButton.x = width / 2;
		this.playButton.y = height / 2;
	}

	public async onClick() {
		console.log("clicked from menu");
		ConnectMatchmakingSocket();
	}
}

function ConnectMatchmakingSocket() {
	const url = "/api/game/match_making/ticket/?auth_token=" + GetJwtToken();
	const socket = new WebSocket(url, []);

	socket.onopen = function(_e: Event) {
		alert("[open] Connection established");
		console.log(_e);
		socket.send("My name is John");
	};

	socket.onmessage = (event) => {
		socket.close(1000, "Normal Closure");

		if (event.data == "found") {
			GetEngine().navigation.showScreen(GameSessionsScreen);
		} else {
			alert("recieved data is " + event.data);
		}
	};

	socket.onclose = function(event) {
		if (event.wasClean) {
			alert(
				`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`,
			);
		} else {
			alert("[close] Connection died");
		}
	};
}
