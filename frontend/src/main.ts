import { gsap } from "gsap";
import { PixiPlugin } from "gsap/PixiPlugin";
import * as PIXI from "pixi.js";

import { setEngine } from "./app/getEngine";
import { LoadScreen } from "./app/screens/LoadScreen";
import { userSettings } from "./app/utils/userSettings";
import { GetJwtToken, ValidateDymmy } from "./auth/jwt";
import { CreationEngine } from "./engine/engine";

/**
 * Importing these modules will automatically register there plugins with the engine.
 */
import "@pixi/sound";
import { MainScreen } from "./app/screens/main/MainScreen";
// import WebApp from "@twa-dev/sdk";
// import { ValidateInitdata } from "./auth/jwt";
// import "@esotericsoftware/spine-pixi-v8";

// WebApp.ready();

async function ensureJwtToken() {
  let token = null;
  while (token == null) {
    const userID = prompt("What is your user ID?", "1");
    if (userID === null) break; // user canceled

    await ValidateDymmy(Number(userID));
    token = GetJwtToken();
  }
}

ensureJwtToken();
// Create a new creation engine instance
const engine = new CreationEngine();
setEngine(engine);

declare global {
  var __PIXI_APP__: CreationEngine;
}

globalThis.__PIXI_APP__ = engine;

(async () => {
  // Initialize the creation engine instance
  await engine.init({
    background: "#1E1E1E",
    resizeOptions: { minWidth: 768, minHeight: 1024, letterbox: false },
  });
  // .addEventListener("contextmenu", (e) => e.preventDefault());

  gsap.registerPlugin(PixiPlugin);
  PixiPlugin.registerPIXI(PIXI);
  // Initialize the user settings
  userSettings.init();

  // Show the load screen
  await engine.navigation.showScreen(LoadScreen);
  // Show the main screen once the load screen is dismissed
  // await engine.navigation.showScreen(XoScreen);
  await engine.navigation.showScreen(MainScreen);
})();
