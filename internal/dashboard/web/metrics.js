"use strict";
// Per-project metrics block of the dashboard card: two SVG bar charts, summary text and the on-demand
// memory/risks <details>. Pure DOM, no library, no remote request; text only through textContent.
(function () {
  var NS = "http://www.w3.org/2000/svg";
  var cache = {}, inflight = {}, openState = {};

  function el(tag, cls, text) {
    var e = document.createElement(tag);
    if (cls) e.className = cls;
    if (text != null) e.textContent = text;
    return e;
  }
  function days(v) { return v === 0 ? "< 1 dia" : String(v) + (v === 1 ? " dia" : " dias"); }
  function kb(bytes) { return String(Math.round(bytes / 1024 * 10) / 10); }

  function chart(title, items, unit) {
    var fig = el("figure", "chart"), cap = el("figcaption");
    fig.style.margin = "8px 0";
    if (!items.length) {
      cap.textContent = title + ": sem dado";
      fig.appendChild(cap);
      return fig;
    }
    cap.textContent = title + ": " + items.map(function (i) { return i.id + " " + String(i.value); }).join(", ");
    var max = Math.max.apply(null, items.map(function (i) { return i.value; })) || 1;
    var w = items.length * 64 + 8, h = 96, top = 14, base = 76;
    var svg = document.createElementNS(NS, "svg");
    svg.setAttribute("viewBox", "0 0 " + w + " " + h);
    svg.setAttribute("role", "img");
    svg.setAttribute("aria-label", title);
    svg.style.width = (w * 1.5) + "px"; // one fixed scale for both charts: same bar and font size
    svg.style.maxWidth = "100%";
    svg.style.height = "auto";
    items.forEach(function (it, i) {
      var x = 8 + i * 64, bh = Math.max(it.value / max * (base - top), it.value > 0 ? 2 : 1);
      var r = document.createElementNS(NS, "rect");
      r.setAttribute("x", x); r.setAttribute("y", base - bh);
      r.setAttribute("width", 48); r.setAttribute("height", bh);
      r.setAttribute("fill", "var(--acc)");
      r.setAttribute("data-prd", it.id);
      r.setAttribute("data-value", String(it.value));
      svg.appendChild(r);
      [[String(it.value), base - bh - 3], [it.id.replace(/^PRD-/, ""), base + 12]].forEach(function (l) {
        var t = document.createElementNS(NS, "text");
        t.setAttribute("x", x + 24); t.setAttribute("y", l[1]);
        t.setAttribute("text-anchor", "middle"); t.setAttribute("font-size", "10");
        t.setAttribute("fill", "currentColor");
        t.textContent = l[0];
        svg.appendChild(t);
      });
    });
    var note = el("div", "mut", unit);
    fig.appendChild(cap); fig.appendChild(svg); fig.appendChild(note);
    return fig;
  }

  function lines(m) {
    var out = [], prds = m.prds || [];
    out.push("incidentes: " + (m.incidents == null ? "ausente ou ilegível (ver avisos)" : String(m.incidents)));
    var mem = m.memory || {};
    if (mem.missing) out.push("memória: ausente");
    else if (mem.lines == null) out.push("memória: ilegível (ver avisos)");
    else out.push("memória: " + String(mem.lines) + " linhas · " + kb(mem.bytes) + " KB");
    var ep = m.epochal || {};
    if (ep.batches == null) out.push("lotes: ilegível (ver avisos)");
    else if (ep.batches === 0 || !ep.last_at) out.push("última consolidação: nunca (" + ep.batches + " lotes)");
    else out.push("última consolidação: " + ep.last_at.slice(0, 10) + " (" + ep.batches + " lotes)");
    var nodate = prds.filter(function (p) { return p.status === "delivered" && p.days == null; }).length;
    out.push("média por PRD: " + (m.prd_days_avg == null ? "sem dado" : days(m.prd_days_avg) + " (n = " + m.prd_days_n + ")") +
      (nodate ? " · " + nodate + " sem data" : ""));
    prds.forEach(function (p) {
      if (p.status === "delivered")
        out.push(p.id + ": " + (p.started || "?") + " → " + (p.delivered || "?") + (p.days == null ? "" : " (" + days(p.days) + ")"));
    });
    prds.forEach(function (p) {
      if (!p.tickets) return;
      out.push(p.id + " · média por ticket: " + (p.ticket_days_avg == null ? "sem dado" : days(p.ticket_days_avg) + " (n = " + p.ticket_days_n + ")") +
        (p.ticket_undated ? " · " + p.ticket_undated + " sem data" : ""));
    });
    return out;
  }

  function load(k, name, file, pre) {
    if (cache[k] != null) { pre.textContent = cache[k]; return; }
    pre.textContent = "carregando...";
    if (!inflight[k]) {
      inflight[k] = fetch("api/memory?project=" + encodeURIComponent(name) + "&file=" + file, { cache: "no-store" })
        .then(function (r) { return r.ok ? r.text().then(function (t) { return [true, t]; }) : [false, "não disponível (HTTP " + r.status + ")"]; })
        .catch(function () { return [false, "não disponível (sem conexão)"]; })
        .then(function (a) { if (a[0]) cache[k] = a[1]; delete inflight[k]; return a[1]; }); // errors are not cached: retry on next open
    }
    inflight[k].then(function (t) { pre.textContent = t; });
  }

  function details(p, file, label) {
    var k = p.path + "#" + file, d = el("details"), pre = el("pre");
    d.dataset.k = k;
    d.appendChild(el("summary", null, label));
    pre.style.cssText = "max-height:18rem;overflow:auto;white-space:pre-wrap;word-break:break-word;margin:4px 0;font-size:.8rem";
    d.appendChild(pre);
    d.addEventListener("toggle", function () {
      openState[k] = d.open;
      if (d.open) load(k, p.name, file, pre);
    });
    if (openState[k]) { d.open = true; load(k, p.name, file, pre); } // rebuilt list: reuse cache, no new request
    return d;
  }

  window.dhMetrics = function (c, p) {
    var m = p.metrics;
    if (!m) return;
    var prds = m.prds || [];
    var del = prds.filter(function (x) { return x.status === "delivered" && x.days != null; })
      .map(function (x) { return { id: x.id, value: x.days }; });
    var nodate = prds.filter(function (x) { return x.status === "delivered" && x.days == null; }).map(function (x) { return x.id; });
    var f1 = chart("dias por PRD entregue", del, "dias do início à entrega");
    if (nodate.length) f1.querySelector("figcaption").textContent += " · sem data: " + nodate.join(", ");
    c.appendChild(f1);
    c.appendChild(chart("tickets por PRD aberto", prds.filter(function (x) { return x.status === "open"; })
      .map(function (x) { return { id: x.id, value: x.tickets }; }), "tickets"));
    var ul = el("ul", "mut");
    lines(m).forEach(function (t) { ul.appendChild(el("li", null, t)); });
    c.appendChild(ul);
    c.appendChild(details(p, "memory", "memória"));
    c.appendChild(details(p, "risks", m.incidents == null ? "incidentes" : "incidentes (" + m.incidents + ")"));
  };
})();
