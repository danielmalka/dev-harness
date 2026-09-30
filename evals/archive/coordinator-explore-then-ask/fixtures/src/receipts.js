// Minimal receipts service.
//
// GET  /receipts/:id      -> return one receipt as JSON
// POST /receipts          -> create a receipt from a support agent's note
//
// No export or email endpoint exists yet.

function getReceipt(id) {
  throw new Error("not implemented");
}

function createReceipt(note) {
  throw new Error("not implemented");
}

module.exports = { getReceipt, createReceipt };
