import App from './App.svelte';

import * as Wails from '@wailsapp/runtime';

let app;

Wails.Init(() => {
  app = new App({
    target: document.getElementById('app')
  });

  window.addEventListener('gamepadconnected', function(e) {
    document.body.innerHTML = 'Gamepad connected at index '
      + e.gamepad.index
      + ': '
      + e.gamepad.id
      + '. '
      + e.gamepad.buttons.length
      + ' buttons, '
      + e.gamepad.axes.length
      + ' axes'
  });
});

export default app;
