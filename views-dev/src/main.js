import Alpine from 'https://cdn.jsdelivr.net/npm/alpinejs@3.12.0/dist/module.esm.min.js'
import { search } from './components/search'

window.Alpine = Alpine

Alpine.data('search', search)

Alpine.start()