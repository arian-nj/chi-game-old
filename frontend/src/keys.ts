import type { Transport } from "@connectrpc/connect";
import type { InjectionKey } from "vue";

const transportKey: InjectionKey<Transport> & symbol = Symbol();
