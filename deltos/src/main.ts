import { decayFactor, weeksBetween, WEEK_MS } from "./engine/decay";
import { netCredits } from "./engine/formula";
import { Ledger } from "./engine/ledger";
import type { NewContribution } from "./engine/types";

const app = document.getElementById("app")!;
const ledger = new Ledger();
let essentialsOptedOut = false;

// Seed with a couple of example contributions so the balance and decay
// aren't zero on first load — same reasoning as toolshed/meshfeed's demo
// seed data: an empty state proves nothing.
const seedNow = Date.now();
ledger.record(
  { description: "Repaired the shared water pump", category: "repairs", laborHours: 6, materialsKg: 2, energyKwh: 1, co2Kg: 0.5 },
  seedNow - 3 * WEEK_MS
);
ledger.record(
  { description: "Delivered firewood by van", category: "transfers", laborHours: 2, materialsKg: 0, energyKwh: 8, co2Kg: 4 },
  seedNow - WEEK_MS
);

function render() {
  const now = Date.now();
  const balance = ledger.currentBalance(now, essentialsOptedOut);
  const decayOnly4wk = ledger.projectDecayOnly(now, 4, essentialsOptedOut);
  const recentAvgRate = averageWeeklyRate();
  const withRate4wk = ledger.projectWithRate(now, 4, recentAvgRate, essentialsOptedOut);

  app.innerHTML = "";

  const balanceEl = document.createElement("div");
  balanceEl.className = `balance ${balance >= 0 ? "positive" : "negative"}`;
  balanceEl.textContent = balance.toFixed(2);
  app.appendChild(balanceEl);

  const essentialsLabel = document.createElement("label");
  const essentialsCheckbox = document.createElement("input");
  essentialsCheckbox.type = "checkbox";
  essentialsCheckbox.checked = essentialsOptedOut;
  essentialsCheckbox.addEventListener("change", () => {
    essentialsOptedOut = essentialsCheckbox.checked;
    render();
  });
  essentialsLabel.appendChild(essentialsCheckbox);
  essentialsLabel.appendChild(document.createTextNode(" Opted out of Community Essentials (steeper 15% weekly decay)"));
  app.appendChild(essentialsLabel);

  const projection = document.createElement("div");
  projection.className = "projection";
  projection.innerHTML = `
    <div><strong>In 4 weeks, if nothing new is contributed:</strong> ${decayOnly4wk.toFixed(2)}</div>
    <div><strong>In 4 weeks, at your recent average rate (${recentAvgRate.toFixed(2)}/week):</strong> ${withRate4wk.toFixed(2)}</div>
  `;
  app.appendChild(projection);

  app.appendChild(buildForm());

  const historyEl = document.createElement("div");
  historyEl.innerHTML = "<h2>Contributions</h2>";
  for (const c of [...ledger.history()].reverse()) {
    const weeksAgo = weeksBetween(c.timestamp, now);
    const row = document.createElement("div");
    row.className = "item";
    const raw = netCredits(c);
    const decayed = raw * decayFactor(weeksAgo, essentialsOptedOut);
    row.innerHTML = `
      <div>
        <div>${c.description} <span class="meta">(${c.category})</span></div>
        <div class="meta">${weeksAgo.toFixed(1)} weeks ago — raw ${raw.toFixed(2)}</div>
      </div>
      <div>${decayed.toFixed(2)}</div>
    `;
    historyEl.appendChild(row);
  }
  app.appendChild(historyEl);
}

function averageWeeklyRate(): number {
  const history = ledger.history();
  if (history.length === 0) return 0;
  const totalRaw = history.reduce((sum, c) => sum + netCredits(c), 0);
  const oldest = Math.min(...history.map((c) => c.timestamp));
  const weeksSpan = Math.max(1, (Date.now() - oldest) / WEEK_MS);
  return totalRaw / weeksSpan;
}

function buildForm(): HTMLFormElement {
  const form = document.createElement("form");
  form.innerHTML = `
    <fieldset>
      <legend>Add a contribution</legend>
      <label>Description <input name="description" type="text" required /></label>
      <label>Category
        <select name="category">
          <option value="labor">Labor</option>
          <option value="materials">Materials</option>
          <option value="repairs">Repairs</option>
          <option value="transfers">Transfers</option>
          <option value="services">Services</option>
        </select>
      </label>
      <label>Labor hours <input name="laborHours" type="number" step="0.1" value="0" /></label>
      <label>Materials (kg) <input name="materialsKg" type="number" step="0.1" value="0" /></label>
      <label>Energy (kWh) <input name="energyKwh" type="number" step="0.1" value="0" /></label>
      <label>CO2 (kg) <input name="co2Kg" type="number" step="0.1" value="0" /></label>
      <button type="submit">Add</button>
    </fieldset>
  `;
  form.addEventListener("submit", (e) => {
    e.preventDefault();
    const data = new FormData(form);
    const input: NewContribution = {
      description: String(data.get("description")),
      category: String(data.get("category")),
      laborHours: Number(data.get("laborHours")),
      materialsKg: Number(data.get("materialsKg")),
      energyKwh: Number(data.get("energyKwh")),
      co2Kg: Number(data.get("co2Kg")),
    };
    ledger.record(input);
    render();
  });
  return form;
}

render();
