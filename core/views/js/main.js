import Alpine from "https://cdn.jsdelivr.net/npm/alpinejs@3.12.0/dist/module.esm.js";
import { search } from "./components/search.js";

window.Alpine = Alpine;

Alpine.data('search', search)

Alpine.start();
