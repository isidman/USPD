import { checkOut, createResource, listResources, returnLoan } from "./api";
import { describeAvailability, kindLabel } from "./format";
import type { ResourceView } from "./types";

const app = document.getElementById("app")!;

// Known limitation, kept deliberately rather than hidden: the API doesn't
// expose "which loan is currently open on this resource" to the list
// endpoint, so a Return button only appears for items checked out in this
// browser tab's own session. Extending GET /resources to include the open
// loan id (if any) is the natural next step — see BLUEPRINT.md.
const loanIdByResource = new Map<string, string>();

async function refresh() {
  const resources = await listResources();
  render(resources);
}

function render(resources: ResourceView[]) {
  app.innerHTML = "";

  const form = document.createElement("form");
  form.innerHTML = `
    <h2>Add a resource</h2>
    <select name="kind">
      <option value="tool">Tool</option>
      <option value="hardware">Hardware</option>
      <option value="software">Software</option>
    </select>
    <input name="name" placeholder="Name" required />
    <button type="submit">Add</button>
  `;
  form.addEventListener("submit", async (e) => {
    e.preventDefault();
    const data = new FormData(form);
    await createResource({
      kind: data.get("kind") as ResourceView["kind"],
      name: String(data.get("name")),
    });
    form.reset();
    await refresh();
  });
  app.appendChild(form);

  const list = document.createElement("div");
  list.innerHTML = "<h2>Resources</h2>";
  for (const resource of resources) {
    const card = document.createElement("div");
    card.className = "resource";
    const statusClass = resource.available ? "available" : "unavailable";
    card.innerHTML = `
      <strong>${kindLabel(resource.kind)}</strong>: ${resource.name}
      — <span class="${statusClass}">${describeAvailability(resource)}</span>
    `;

    if (resource.available) {
      const button = document.createElement("button");
      button.textContent = "Check out";
      button.addEventListener("click", async () => {
        const borrower = prompt("Your name?");
        if (!borrower) return;
        const loan = await checkOut(resource.id, borrower);
        loanIdByResource.set(resource.id, loan.id);
        await refresh();
      });
      card.appendChild(button);
    } else if (loanIdByResource.has(resource.id)) {
      const button = document.createElement("button");
      button.textContent = "Return";
      button.addEventListener("click", async () => {
        const loanId = loanIdByResource.get(resource.id);
        if (!loanId) return;
        await returnLoan(loanId);
        loanIdByResource.delete(resource.id);
        await refresh();
      });
      card.appendChild(button);
    }

    list.appendChild(card);
  }
  app.appendChild(list);
}

refresh().catch((err) => {
  app.textContent = `Failed to load resources: ${err.message}`;
});
