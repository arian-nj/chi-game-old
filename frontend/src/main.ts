import { GetEngine, setEngine } from "./app/getEngine";
import { LoadScreen } from "./app/screens/LoadScreen";
import { XoScreen } from "./app/screens/xo/XOScreen";
import { userSettings } from "./app/utils/userSettings";
import { CreationEngine } from "./engine/engine";
/**
 * Importing these modules will automatically register there plugins with the engine.
 */
import "@pixi/sound";
// import WebApp from "@twa-dev/sdk";
// import { ValidateInitdata } from "./auth/jwt";
// import "@esotericsoftware/spine-pixi-v8";
//
// WebApp.ready();
// (async () => {
//   await ValidateInitdata(WebApp.initData);
// })();
//
// Create a new creation engine instance
const engine = new CreationEngine();
setEngine(engine);

globalThis.__PIXI_APP__ = engine;
(async () => {
  // Initialize the creation engine instance
  await engine.init({
    background: "#1E1E1E",
    resizeOptions: { minWidth: 768, minHeight: 1024, letterbox: false },
  });

  // Initialize the user settings
  userSettings.init();

  // Show the load screen
  await engine.navigation.showScreen(LoadScreen);
  // Show the main screen once the load screen is dismissed
  await engine.navigation.showScreen(XoScreen);
  // await engine.navigation.showScreen(MainScreen);
})();
