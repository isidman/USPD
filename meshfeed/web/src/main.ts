import { listFeeds, listItems } from "./api";
import { formatTimestamp, shortFeedId } from "./format";

const app = document.getElementById("app")!;

async function render() {
  const feeds = await listFeeds();
  app.innerHTML = "";

  if (feeds.length === 0) {
    app.textContent = "No feeds held yet — nothing published or synced so far.";
    return;
  }

  for (const feed of feeds) {
    const card = document.createElement("div");
    card.className = "feed";
    card.innerHTML = `<strong>${shortFeedId(feed.id)}</strong> <span class="meta">(latest seq ${feed.latest_seq})</span>`;

    const items = await listItems(feed.id);
    for (const item of items.slice().reverse()) {
      const row = document.createElement("div");
      row.className = "item";
      row.innerHTML = `<div>${item.content}</div><div class="meta">#${item.seq} — ${formatTimestamp(item.timestamp)}</div>`;
      card.appendChild(row);
    }

    app.appendChild(card);
  }
}

render().catch((err) => {
  app.textContent = `Failed to load feeds: ${err.message}`;
});
