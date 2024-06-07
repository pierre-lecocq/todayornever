// File: todayornever.js
// Creation: Fri Jun  7 09:18:20 2024
// Time-stamp: <2024-06-07 14:52:12>
// Copyright (C): 2024 Pierre Lecocq

htmx.on('htmx:oobAfterSwap', evt => {
  if (evt.detail.target && evt.detail.target.id === 'feedback') {
    evt.detail.target.style['display'] = 'block';
    setTimeout(() => {
      evt.detail.target.style['display'] = 'none';
    }, 5000)
  }
})
