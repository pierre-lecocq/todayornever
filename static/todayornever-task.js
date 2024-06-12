// File: todayornever-task.js
// Creation: Tue Jun 11 11:21:41 2024
// Time-stamp: <2024-06-11 18:11:45>
// Copyright (C): 2024 Pierre Lecocq

// -----------------------------------------------------------------------------
// Edit

var dialogEditElements = document.getElementsByClassName('dialog-edit');

for (let i = 0; i < dialogEditElements.length; i++) {
  let display = dialogEditElements[i].getElementsByTagName('span')[0];
  let dialog = dialogEditElements[i].getElementsByTagName('dialog')[0];

  display.addEventListener('dblclick', function (e) {
    if (!dialog.open) {
      dialog.showModal();
    }
  });
}

// -----------------------------------------------------------------------------
// Sort

var draggableElements = document.getElementsByClassName('grip');

function dragEvent(event) {
	event.dataTransfer.setData("text/plain", this.dataset.id);
}

function dragOverEvent(event) {
	event.preventDefault();
}

function dropEvent(event) {
	event.preventDefault();
	let id1 = event.dataTransfer.getData("text/plain");
  let id2 = this.dataset.id;
  this.classList.remove('drop-hover');

  if (id1 && id2 && id1 !== id2) {
    htmx.ajax('POST', '/reorder', { swap: 'none', values: { id1, id2 } });
  }
}

function dragEnterEvent(event) {
  this.classList.add('drop-hover');
}

function dragLeaveEvent(event) {
  this.classList.remove('drop-hover');
}

for (let i = 0; i < draggableElements.length; i++) {
  let item = draggableElements[i];
  item.parentNode.addEventListener('dragover', dragOverEvent);
  item.parentNode.addEventListener('dragstart', dragEvent);
  item.parentNode.addEventListener('drop', dropEvent);
  item.parentNode.addEventListener('dragenter', dragEnterEvent);
  item.parentNode.addEventListener('dragleave', dragLeaveEvent);
}
